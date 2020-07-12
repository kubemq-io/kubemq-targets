package events_store

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
		kubemq.WithCheckConnection(true),
	)

	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	eventStoreMetadata, err := parseMetadata(request.Metadata, c.opts)
	if err != nil {
		return nil, err
	}
	result, err := c.client.ES().
		SetId(eventStoreMetadata.id).
		SetChannel(eventStoreMetadata.channel).
		SetMetadata(eventStoreMetadata.metadata).
		SetBody(request.Data).
		Send(ctx)
	if err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, fmt.Errorf(result.Err.Error())
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("id", result.Id), nil
}
