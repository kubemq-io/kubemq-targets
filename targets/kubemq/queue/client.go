package queue

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
	"github.com/kubemq-io/kubemq-go"
)

type Client struct {
	name   string
	opts   options
	client *kubemq.Client
}

func New() *Client {
	return &Client{}

}
func (c *Client) Name() string {
	return c.name
}
func (c *Client) Init(ctx context.Context, cfg config.Metadata) error {
	c.name = cfg.Name
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	c.client, err = kubemq.NewClient(ctx,
		kubemq.WithAddress(c.opts.host, c.opts.port),
		kubemq.WithClientId(c.opts.clientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithAuthToken(c.opts.authToken),
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	queueMetadata, err := parseMetadata(request.Metadata, c.opts)
	if err != nil {
		return nil, err
	}
	queueMessage := c.client.NewQueueMessage().
		SetId(queueMetadata.id).
		SetChannel(queueMetadata.channel).
		SetMetadata(queueMetadata.metadata).
		SetBody(request.Data).
		SetPolicyDelaySeconds(queueMetadata.delaySeconds).
		SetPolicyExpirationSeconds(queueMetadata.expirationSeconds).
		SetPolicyMaxReceiveCount(queueMetadata.maxReceiveCount).
		SetPolicyMaxReceiveQueue(queueMetadata.deadLetterQueue)
	result, err := queueMessage.Send(ctx)
	if err != nil {
		return nil, err
	}
	if result.IsError {
		return nil, fmt.Errorf(result.Error)
	}
	return types.NewResponse().
			SetMetadataKeyValue("id", queueMetadata.id).
			SetMetadataKeyValue("result", "ok"),
		nil
}
