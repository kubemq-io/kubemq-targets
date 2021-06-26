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
	"github.com/kubemq-io/kubemq-go/queues_stream"
	"time"
)

var (
	errInvalidTarget = errors.New("invalid controller received, cannot be null")
)

type Client struct {
	opts      options
	client    *queues_stream.QueuesStreamClient
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
	c.client, err = queues_stream.NewQueuesStreamClient(ctx,
		queues_stream.WithAddress(c.opts.host, c.opts.port),
		queues_stream.WithClientId(c.opts.clientId),
		queues_stream.WithCheckConnection(true),
		queues_stream.WithAuthToken(c.opts.authToken))
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
		err := c.processQueueMessage(ctx)
		if err != nil {
			c.log.Error(err.Error())
			time.Sleep(time.Second)
		}
		select {
		case <-ctx.Done():
			return
		default:

		}
	}
}
func (c *Client) processQueueMessage(ctx context.Context)  error {
	pr := queues_stream.NewPollRequest().
		SetChannel(c.opts.channel).
		SetMaxItems(c.opts.batchSize).
		SetWaitTimeout(c.opts.waitTimeout).
		SetAutoAck(false)
	pollResp, err := c.client.Poll(ctx, pr)
	if err != nil {
		return err
	}
	if !pollResp.HasMessages() {
		return nil
	}

	for _, message := range pollResp.Messages {
		req, err := types.ParseRequest(message.Body)
		if err != nil {
			_ = message.Ack()
			return fmt.Errorf("invalid request format, %w", err)
		}
		resp, err := c.target.Do(ctx, req)
		if err != nil {
			if message.Policy.MaxReceiveCount != message.Attributes.ReceiveCount {
				return message.NAck()
			}
			if c.opts.responseChannel != "" {
				errResp:=types.NewResponse().SetError(err)
				_, errSend := c.client.Send(ctx, errResp.ToQueueStreamMessage().SetChannel(c.opts.responseChannel))
				if errSend != nil {
					c.log.Errorf("error sending response to a queue, %s", errSend.Error())
				}
			}
			return nil
		}
		if c.opts.resend != "" {
			_ = message.ReQueue(c.opts.resend)
		} else {
			_ = message.Ack()
		}
		if resp != nil {
			if c.opts.responseChannel != "" {
				_, errSend := c.client.Send(ctx, resp.ToQueueStreamMessage().SetChannel(c.opts.responseChannel))
				if errSend != nil {
					c.log.Errorf("error sending response to a queue, %s", errSend.Error())
				}
			}
		}
	}
	return nil
}

func (c *Client) Stop() error {
	c.isStopped = true
	_ = c.client.Close()
	return nil
}
