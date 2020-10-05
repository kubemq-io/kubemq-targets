package servicebus

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Azure/azure-service-bus-go"
	"github.com/kubemq-hub/builder/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type Client struct {
	name   string
	opts   options
	client *servicebus.Queue
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
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(c.opts.connectionString))
	if err != nil {
		return err
	}
	c.client, err = ns.NewQueue(c.opts.queueName)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "send":
		return c.send(ctx, meta, req.Data)
	case "send_batch":
		return c.sendBatch(ctx, meta, req.Data)
	}
	return nil, errors.New("invalid method type")
}

func (c *Client) send(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	m := servicebus.NewMessage(data)
	if data == nil {
		return nil, errors.New("missing data")
	}
	if meta.label != "" {
		m.Label = meta.label
	}
	if meta.contentType != "" {
		m.ContentType = meta.contentType
	}
	if meta.timeToLive > 0 {
		m.TTL = &meta.timeToLive
	}
	err := c.client.Send(ctx, m)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) sendBatch(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	var messages []string
	err := json.Unmarshal(data, &messages)
	if err != nil {
		return nil, err
	}
	var sm []*servicebus.Message
	for _, m := range messages {
		message := servicebus.NewMessageFromString(m)
		sm = append(sm, message)
	}

	for _, m := range sm {
		if meta.timeToLive > 0 {
			m.TTL = &meta.timeToLive
		}
		if meta.label != "" {
			m.Label = meta.label
		}
		if meta.contentType != "" {
			m.ContentType = meta.contentType
		}
	}

	err = c.client.SendBatch(ctx, servicebus.NewMessageBatchIterator(meta.maxBatchSize, sm...))
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.client.Close(ctx)
}
