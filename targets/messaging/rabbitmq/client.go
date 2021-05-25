package rabbitmq

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/streadway/amqp"
	"sync"
)

type Client struct {
	sync.Mutex
	name        string
	opts        options
	channel     *amqp.Channel
	conn        *amqp.Connection
	isConnected bool
}

func New() *Client {
	return &Client{
		name:    "",
		opts:    options{},
		channel: nil,
	}
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
	if err := c.connect(); err != nil {
		return err
	}
	return nil
}
func (c *Client) connect() error {
	var err error
	c.conn, err = amqp.Dial(c.opts.url)
	if err != nil {
		return fmt.Errorf("error dialing rabbitmq, %w", err)
	}
	c.channel, err = c.conn.Channel()
	if err != nil {
		_ = c.conn.Close()
		return fmt.Errorf("error getting rabbitmq channel, %w", err)
	}
	c.isConnected = true
	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, ok := c.opts.defaultMetadata()
	if !ok {
		var err error
		meta, err = parseMetadata(req.Metadata)
		if err != nil {
			return nil, err
		}
	}
	return c.Publish(ctx, meta, req.Data)
}

func (c *Client) Publish(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	c.Lock()
	defer c.Unlock()
	if !c.isConnected {
		if err := c.connect(); err != nil {
			return nil, err
		}
	}
	msg := meta.amqpMessage(data)
	err := c.channel.Publish(meta.exchange, meta.queue, meta.mandatory, meta.immediate, msg)
	if err != nil {
		c.isConnected = false
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) Stop() error {
	if c.channel != nil {
		return c.channel.Close()
	}
	return nil
}
