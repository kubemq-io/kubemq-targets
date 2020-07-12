package events

import (
	"context"
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
	eventMetadata, err := parseMetadata(request.Metadata, c.opts)
	if err != nil {
		return nil, err
	}
	err = c.client.E().
		SetId(eventMetadata.id).
		SetChannel(eventMetadata.channel).
		SetMetadata(eventMetadata.metadata).
		SetBody(request.Data).
		Send(ctx)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetMetadataKeyValue("id", eventMetadata.id),
		nil
}
