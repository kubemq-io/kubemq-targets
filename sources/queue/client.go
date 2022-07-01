package queue

import (
	"context"
	"errors"
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-go/queues_stream"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/middleware"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
	"time"
)

var (
	errInvalidTarget = errors.New("invalid controller received, cannot be null")
)

type Client struct {
	opts      options
	log       *logger.Logger
	target    middleware.Middleware
	isStopped bool
}

func (c *Client) getQueuesClient(ctx context.Context, id int) (*queues_stream.QueuesStreamClient, error) {
	return queues_stream.NewQueuesStreamClient(ctx,
		queues_stream.WithAddress(c.opts.host, c.opts.port),
		queues_stream.WithClientId(c.opts.clientId),
		queues_stream.WithCheckConnection(true),
		queues_stream.WithAutoReconnect(true),
		queues_stream.WithAuthToken(c.opts.authToken),
		queues_stream.WithConnectionNotificationFunc(
			func(msg string) {
				c.log.Infof(fmt.Sprintf("connection: %d, %s", id, msg))
			}),
	)

}
func New() *Client {
	return &Client{}

}
func (c *Client) Connector() *common.Connector {
	return Connector()
}
func (c *Client) onError(err error) {
	c.log.Error(err.Error())
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

	return nil
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {
	if target == nil {
		return errInvalidTarget
	} else {
		c.target = target
	}
	for i := 0; i < c.opts.sources; i++ {
		client, err := c.getQueuesClient(ctx, i+1)
		if err != nil {
			return err
		}
		go c.run(ctx, client)
	}
	return nil
}

func (c *Client) run(ctx context.Context, client *queues_stream.QueuesStreamClient) {
	defer func() {
		_ = client.Close()
	}()
	for {
		if c.isStopped {
			return
		}
		err := c.processQueueMessage(ctx, client)
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
func (c *Client) processQueueMessage(ctx context.Context, client *queues_stream.QueuesStreamClient) error {
	pr := queues_stream.NewPollRequest().
		SetChannel(c.opts.channel).
		SetMaxItems(c.opts.batchSize).
		SetWaitTimeout(c.opts.waitTimeout * 1000).
		SetAutoAck(true).
		SetOnErrorFunc(c.onError)
	pollResp, err := client.Poll(ctx, pr)
	if err != nil {
		return err
	}
	if !pollResp.HasMessages() {
		return nil
	}
	for _, message := range pollResp.Messages {
		var req *types.Request
		var err error
		if c.opts.doNotParsePayload {
			req = types.NewRequest().SetData(message.Body)
		} else {
			req, err = types.ParseRequest(message.Body)
			if err != nil {
				return fmt.Errorf("invalid request format, %w", err)
			}
		}
		resp, err := c.target.Do(ctx, req)
		if err != nil {
			if c.opts.responseChannel != "" {
				errResp := types.NewResponse().SetError(err)
				_, errSend := client.Send(ctx, errResp.ToQueueStreamMessage().SetChannel(c.opts.responseChannel))
				if errSend != nil {
					c.log.Errorf("error sending response to a queue, %s", errSend.Error())
				}
			}
		}
		if resp != nil {
			if c.opts.responseChannel != "" {
				_, errSend := client.Send(ctx, resp.ToQueueStreamMessage().SetChannel(c.opts.responseChannel))
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
	return nil
}
