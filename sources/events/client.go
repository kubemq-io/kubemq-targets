package events

import (
	"context"
	"errors"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/middleware"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
)

var (
	errInvalidTarget = errors.New("invalid target received, cannot be nil")
)

type Client struct {
	name   string
	opts   options
	client *kubemq.Client
	log    *logger.Logger
	target middleware.Middleware
}

func New() *Client {
	return &Client{}

}
func (c *Client) Name() string {
	return c.name
}
func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	c.name = cfg.Name
	c.log = logger.NewLogger(cfg.Name)
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
		kubemq.WithCheckConnection(false),
		kubemq.WithMaxReconnects(c.opts.maxReconnects),
		kubemq.WithAutoReconnect(c.opts.autoReconnect),
		kubemq.WithReconnectInterval(c.opts.reconnectIntervalSeconds))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {
	if target == nil {
		return errInvalidTarget
	} else {
		c.target = target
	}
	group := nuid.Next()
	if c.opts.group != "" {
		group = c.opts.group
	}

	errCh := make(chan error, 1)
	eventsCh, err := c.client.SubscribeToEvents(ctx, c.opts.channel, group, errCh)
	if err != nil {
		return fmt.Errorf("error on subscribing to events channel, %w", err)
	}
	go func(ctx context.Context, eventsCh <-chan *kubemq.Event, errCh chan error) {
		c.run(ctx, eventsCh, errCh)
	}(ctx, eventsCh, errCh)

	return nil
}

func (c *Client) run(ctx context.Context, eventsCh <-chan *kubemq.Event, errCh chan error) {
	for {
		select {
		case event := <-eventsCh:
			go func(event *kubemq.Event) {
				resp, err := c.processEvent(ctx, event)
				if err != nil {
					resp = types.NewResponse().SetError(err)
				}
				if c.opts.responseChannel != "" {
					errSend := c.client.SetEvent(resp.ToEvent()).SetChannel(c.opts.responseChannel).Send(ctx)
					if errSend != nil {
						c.log.Errorf("error sending event response %s", errSend.Error())
					}
				}
			}(event)

		case err := <-errCh:
			c.log.Errorf("error received from kuebmq server, %s", err.Error())
			return
		case <-ctx.Done():
			return

		}
	}
}

func (c *Client) processEvent(ctx context.Context, event *kubemq.Event) (*types.Response, error) {
	req, err := types.ParseRequest(event.Body)
	if err != nil {
		return nil, fmt.Errorf("invalid request format, %w", err)
	}
	resp, err := c.target.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *Client) Stop() error {
	return c.client.Close()
}
