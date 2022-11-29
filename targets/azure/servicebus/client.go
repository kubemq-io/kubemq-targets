package servicebus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Azure/azure-service-bus-go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
)

type Client struct {
	log    *logger.Logger
	opts   options
	client *servicebus.Queue
}

func New() *Client {
	return &Client{}
}

func (c *Client) Connector() *common.Connector {
	return Connector()
}

func (c *Client) Init(ctx context.Context, cfg config.Spec, log *logger.Logger) error {
	c.log = log
	if c.log == nil {
		c.log = logger.NewLogger(cfg.Kind)
	}

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
		return fmt.Errorf("error connecting to servicebus at %s: %w", c.opts.connectionString, err)
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

func (c *Client) Stop() error {
	if c.client != nil {
		return c.client.Close(context.Background())
	}
	return nil
}
