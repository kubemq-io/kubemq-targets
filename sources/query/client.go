package query

import (
	"context"
	"errors"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/middleware"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/nats-io/nuid"

	"github.com/kubemq-io/kubemq-go"
	"time"
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
	c.client, _ = kubemq.NewClient(ctx,
		kubemq.WithAddress(c.opts.host, c.opts.port),
		kubemq.WithClientId(c.opts.clientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithAuthToken(c.opts.authToken),
		kubemq.WithMaxReconnects(c.opts.maxReconnects),
		kubemq.WithAutoReconnect(c.opts.autoReconnect),
		kubemq.WithReconnectInterval(c.opts.reconnectIntervalSeconds))
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
		queriesCh, err := c.client.SubscribeToQueries(ctx, c.opts.channel, group, errCh)
		if err != nil {
			return fmt.Errorf("error on subscribing to query channel, %w", err)
		}
		go func(ctx context.Context, queryCh <-chan *kubemq.QueryReceive, errCh chan error) {
			c.run(ctx, queriesCh, errCh)
		}(ctx, queriesCh, errCh)
	}
	return nil
}

func (c *Client) run(ctx context.Context, queryCh <-chan *kubemq.QueryReceive, errCh chan error) {
	for {
		select {
		case query := <-queryCh:

			go func(q *kubemq.QueryReceive) {
				queryResponse := c.client.R().
					SetRequestId(query.Id).
					SetResponseTo(query.ResponseTo)
				resp, err := c.processQuery(ctx, query)
				if err != nil {
					resp = types.NewResponse().SetError(err)
				}
				queryResponse.SetExecutedAt(time.Now()).
					SetBody(resp.MarshalBinary())
				err = queryResponse.Send(ctx)
				if err != nil {
					c.log.Errorf("error sending query response %s", err.Error())
				}
			}(query)

		case err := <-errCh:
			c.log.Errorf("error received from kuebmq server, %s", err.Error())
			return
		case <-ctx.Done():
			return

		}
	}
}

func (c *Client) processQuery(ctx context.Context, query *kubemq.QueryReceive) (*types.Response, error) {
	req, err := types.ParseRequest(query.Body)
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
