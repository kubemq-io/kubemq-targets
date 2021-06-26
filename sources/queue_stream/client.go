package queue_stream

import (
	"context"
	"errors"
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/middleware"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"strings"
	"time"
)

var (
	errInvalidTarget = errors.New("invalid controller received, cannot be null")
)

type Client struct {
	opts      options
	client    *kubemq.Client
	log       *logger.Logger
	target    middleware.Middleware
	isStopped bool
}

func New() *Client {
	return &Client{}

}
func (c *Client) Connector() *common.Connector {
	return Connector()
}
func (c *Client) getKubemqClient(ctx context.Context) (*kubemq.Client, error) {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress(c.opts.host, c.opts.port),
		kubemq.WithClientId(c.opts.clientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithAuthToken(c.opts.authToken))
	if err != nil {
		return nil, err
	}
	return client, nil
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
	c.client, err = kubemq.NewClient(ctx,
		kubemq.WithAddress(c.opts.host, c.opts.port),
		kubemq.WithClientId(c.opts.clientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true),
		kubemq.WithAuthToken(c.opts.authToken))
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
	for i := 0; i < c.opts.sources; i++ {
		go c.run(ctx)
	}
	return nil
}

func (c *Client) run(ctx context.Context) {
	for {
		if c.isStopped {
			return
		}
		resp, err := c.processQueueMessage()
		if err != nil {
			if !strings.Contains(err.Error(), "138") {
				c.log.Error(err.Error())
				time.Sleep(time.Second)
			}
		} else {
			if resp != nil {
				if c.opts.responseChannel != "" {
					_, errSend := c.client.SetQueueMessage(resp.ToQueueMessage()).SetChannel(c.opts.responseChannel).Send(ctx)
					if errSend != nil {
						c.log.Errorf("error sending response to a queue, %s", errSend.Error())
					}
				}
			}
		}
		select {
		case <-ctx.Done():
			return
		default:

		}
	}
}
func (c *Client) processQueueMessage() (*types.Response, error) {
	ctx := context.Background()
	client, err := c.getKubemqClient(ctx)
	if err != nil {
		return nil, err
	}
	stream := client.NewStreamQueueMessage().SetChannel(c.opts.channel)
	defer func() {
		stream.Close()
	}()
	msg, err := stream.Next(ctx, int32(c.opts.visibilityTimeout), int32(c.opts.waitTimeout))
	if err != nil {
		return nil, err
	}

	req, err := types.ParseRequest(msg.Body)
	if err != nil {
		return types.NewResponse().SetError(fmt.Errorf("invalid request format, %w", err)), msg.Ack()
	}
	resp, err := c.target.Do(ctx, req)
	if err != nil {
		if msg.Policy.MaxReceiveCount != msg.Attributes.ReceiveCount {
			return nil, msg.Reject()
		}
		return types.NewResponse().SetError(err), nil
	}
	if c.opts.resend != "" {
		return resp, msg.Resend(c.opts.resend)
	}
	return resp, msg.Ack()
}

func (c *Client) Stop() error {
	c.isStopped = true
	_ = c.client.Close()
	return nil
}
