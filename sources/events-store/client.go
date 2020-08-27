package events_store

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

const (
	defaultHost          = "localhost"
	defaultPort          = 50000
	defaultAutoReconnect = true
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
		kubemq.WithCheckConnection(true),
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
	for i := 0; i < c.opts.concurrency; i++ {
		errCh := make(chan error, 1)
		eventsCh, err := c.client.SubscribeToEventsStore(ctx, c.opts.channel, group, errCh, kubemq.StartFromNewEvents())
		if err != nil {
			return fmt.Errorf("error on subscribing to events channel, %w", err)
		}
		go func(ctx context.Context, eventsCh <-chan *kubemq.EventStoreReceive, errCh chan error) {
			c.run(ctx, eventsCh, errCh)
		}(ctx, eventsCh, errCh)
	}
	return nil
}

func (c *Client) run(ctx context.Context, eventsCh <-chan *kubemq.EventStoreReceive, errCh chan error) {
	for {
		select {
		case event := <-eventsCh:
			go func(event *kubemq.EventStoreReceive) {
				resp, err := c.processEventStore(ctx, event)
				if err != nil {
					resp = types.NewResponse().SetError(err)
				}
				if c.opts.responseChannel != "" {
					sendRes, errSend := c.client.SetEventStore(resp.ToEventStore()).SetChannel(c.opts.responseChannel).Send(ctx)
					if errSend != nil {
						c.log.Errorf("error sending event response %s", errSend.Error())
					} else {
						if !sendRes.Sent {
							c.log.Errorf("error sending event response %s", sendRes.Err)
						}
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

func (c *Client) processEventStore(ctx context.Context, event *kubemq.EventStoreReceive) (*types.Response, error) {
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
