package query

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/middleware"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
	"github.com/kubemq-io/kubemq-targets/types"

	"github.com/kubemq-io/kubemq-go"
)

var errInvalidTarget = errors.New("invalid target received, cannot be nil")

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

func (c *Client) Init(ctx context.Context, cfg config.Spec, bindingName string, log *logger.Logger) error {
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
			clientId = fmt.Sprintf("kubemq-bridges/%s/%s/%d", bindingName, clientId, i)
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
	queriesCh, err := client.SubscribeToQueries(ctx, c.opts.channel, c.opts.group, errCh)
	if err != nil {
		return fmt.Errorf("error on subscribing to query channel, %w", err)
	}
	go func(ctx context.Context, queryCh <-chan *kubemq.QueryReceive, errCh chan error) {
		for {
			select {
			case query := <-queryCh:

				go func(q *kubemq.QueryReceive) {
					queryResponse := client.R().
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
	}(ctx, queriesCh, errCh)
	return nil
}

func (c *Client) processQuery(ctx context.Context, query *kubemq.QueryReceive) (*types.Response, error) {
	var req *types.Request
	var err error
	if c.opts.doNotParsePayload {
		req = types.NewRequest().SetData(query.Body)
	} else {
		req, err = types.ParseRequest(query.Body)
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
