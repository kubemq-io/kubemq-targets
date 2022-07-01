package events_store

import (
	"context"
	"errors"
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/middleware"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"
)

var (
	errInvalidTarget = errors.New("invalid target received, cannot be nil")
)

type Client struct {
	opts    options
	clients []*kubemq.Client
	log     *logger.Logger
	target  middleware.Middleware
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
	for i := 0; i < c.opts.sources; i++ {
		clientId := c.opts.clientId
		if c.opts.sources > 1 {
			clientId = fmt.Sprintf("%s-%d", clientId, i)
		}
		client, err := kubemq.NewClient(ctx,
			kubemq.WithAddress(c.opts.host, c.opts.port),
			kubemq.WithClientId(clientId),
			kubemq.WithTransportType(kubemq.TransportTypeGRPC),
			kubemq.WithCheckConnection(true),
			kubemq.WithAuthToken(c.opts.authToken),
			kubemq.WithMaxReconnects(c.opts.maxReconnects),
			kubemq.WithAutoReconnect(c.opts.autoReconnect),
			kubemq.WithReconnectInterval(c.opts.reconnectIntervalSeconds))
		if err != nil {
			return err
		}
		c.clients = append(c.clients, client)
	}
	return nil
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {
	if target == nil {
		return errInvalidTarget
	} else {
		c.target = target
	}
	if c.opts.sources > 1 && c.opts.group == "" {
		c.opts.group = uuid.New().String()
	}

	for i := 0; i < len(c.clients); i++ {
		err := c.runClient(ctx, c.clients[i])
		if err != nil {
			return fmt.Errorf("error during start of client %d: %s", i, err.Error())
		}
	}
	return nil
}
func (c *Client) runClient(ctx context.Context, client *kubemq.Client) error {
	errCh := make(chan error, 1)
	eventsCh, err := client.SubscribeToEventsStore(ctx, c.opts.channel, c.opts.group, errCh, kubemq.StartFromNewEvents())
	if err != nil {
		return fmt.Errorf("error on subscribing to events channel, %w", err)
	}
	go func(ctx context.Context, eventsCh <-chan *kubemq.EventStoreReceive, errCh chan error) {
		for {
			select {
			case event := <-eventsCh:
				go func(event *kubemq.EventStoreReceive) {
					resp, err := c.processEventStore(ctx, event)
					if err != nil {
						resp = types.NewResponse().SetError(err)
					}
					if c.opts.responseChannel != "" {
						sendRes, errSend := client.SetEventStore(resp.ToEventStore()).SetChannel(c.opts.responseChannel).Send(ctx)
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
	}(ctx, eventsCh, errCh)
	return nil
}

func (c *Client) processEventStore(ctx context.Context, event *kubemq.EventStoreReceive) (*types.Response, error) {
	var req *types.Request
	var err error
	if c.opts.doNotParsePayload {
		req = types.NewRequest().SetData(event.Body)
	} else {
		req, err = types.ParseRequest(event.Body)
		if err != nil {
			return nil, fmt.Errorf("invalid request format, %w", err)
		}
	}
	resp, err := c.target.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *Client) Stop() error {
	for _, client := range c.clients {
		_ = client.Close()
	}
	return nil
}
