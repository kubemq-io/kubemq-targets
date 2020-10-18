package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Client struct {
	name   string
	opts   options
	client *pubsub.Client
}

func New() *Client {
	return &Client{}

}
func (c *Client) Connector() *common.Connector {
	return Connector()
}
func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	c.name = cfg.Name
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	b := []byte(c.opts.credentials)

	client, err := pubsub.NewClient(ctx, c.opts.projectID, option.WithCredentialsJSON(b))
	if err != nil {
		return err
	}
	c.client = client
	return nil
}

func (c *Client) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	eventMetadata, err := parseMetadata(request.Metadata, c.opts)
	if err != nil {
		return nil, err
	}
	t := c.client.Topic(eventMetadata.topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data:       request.Data,
		Attributes: eventMetadata.tags,
	})
	tries := 0
	for tries <= c.opts.retries {
		id, err := result.Get(ctx)
		if err == nil {
			return types.NewResponse().
					SetMetadataKeyValue("event_id", id),
				nil
		}
		if tries >= c.opts.retries {
			return nil, err
		}
		tries++
	}
	return nil, fmt.Errorf("retries must be a zero or greater")
}

func (c *Client) list(ctx context.Context) (*types.Response, error) {
	var topics []string
	it := c.client.Topics(ctx)
	for {
		topic, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic.ID())
	}
	if len(topics) <= 0 {
		return nil, fmt.Errorf("no topics found for this project")
	}
	data, err := json.Marshal(topics)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetData(data).
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) CloseClient() error {
	return c.client.Close()
}
