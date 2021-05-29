package queue

import (
	"context"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/middleware"
	"github.com/kubemq-hub/kubemq-targets/types"
	"github.com/kubemq-io/kubemq-go"
	"time"

	"errors"
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
)

var (
	errInvalidTarget = errors.New("invalid controller received, cannot be null")
)

const (
	retriesInterval = 1 * time.Second
)

type Client struct {
	opts         options
	clients      []*kubemq.Client
	log          *logger.Logger
	target       middleware.Middleware
	isStopped    bool
	requeueCache *requeue
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
	c.requeueCache = newRequeue(c.opts.maxRequeue)
	for i := 0; i < c.opts.sources; i++ {
		clientId := c.opts.clientId
		if c.opts.sources > 1 {
			clientId = fmt.Sprintf("%s-%d", clientId, i)
		}
		client, err := kubemq.NewClient(ctx,
			kubemq.WithAddress(c.opts.host, c.opts.port),
			kubemq.WithClientId(clientId),
			kubemq.WithTransportType(kubemq.TransportTypeGRPC),
			kubemq.WithAuthToken(c.opts.authToken),
			kubemq.WithCheckConnection(true))
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
	for i := 0; i < len(c.clients); i++ {
		go c.run(ctx, c.clients[i])
	}
	return nil
}

func (c *Client) run(ctx context.Context, client *kubemq.Client) {
	for {
		if c.isStopped {
			return
		}
		queueMessages, err := c.getQueueMessages(ctx, client)
		if err != nil {
			c.log.Error(err.Error())
			time.Sleep(retriesInterval)
			continue
		}
		for _, message := range queueMessages {
			resp, err := c.processQueueMessage(ctx, message, client)
			if err != nil {
				c.log.Error(err.Error())
				continue
			}
			if resp != nil && c.opts.responseChannel != "" {
				_, errSend := client.SetQueueMessage(resp.ToQueueMessage()).SetChannel(c.opts.responseChannel).Send(ctx)
				if errSend != nil {
					c.log.Errorf("error sending response to a queue, %s", errSend.Error())
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
func (c *Client) getQueueMessages(ctx context.Context, client *kubemq.Client) ([]*kubemq.QueueMessage, error) {
	receiveResult, err := client.NewReceiveQueueMessagesRequest().
		SetChannel(c.opts.channel).
		SetMaxNumberOfMessages(c.opts.batchSize).
		SetWaitTimeSeconds(c.opts.waitTimeout).
		Send(ctx)
	if err != nil {
		return nil, err
	}
	return receiveResult.Messages, nil
}

func (c *Client) processQueueMessage(ctx context.Context, msg *kubemq.QueueMessage, client *kubemq.Client) (*types.Response, error) {
	req, err := types.ParseRequest(msg.Body)
	if err != nil {
		return types.NewResponse().SetError(fmt.Errorf("invalid request format, %w", err)), nil
	}
	resp, err := c.target.Do(ctx, req)
	if err == nil {
		c.requeueCache.remove(msg.MessageID)
		return resp, nil
	}
	if c.requeueCache.isRequeue(msg.MessageID) {
		_, requeueErr := client.SetQueueMessage(msg).Send(ctx)
		if requeueErr != nil {
			c.requeueCache.remove(msg.MessageID)
			c.log.Errorf("message id %s wasn't requeue due to an error , %s", msg.MessageID, requeueErr.Error())
			return types.NewResponse().SetError(err), nil
		}
		c.log.Infof("message id %s, requeued back to channel", msg.MessageID)
		return nil, nil
	} else {
		return types.NewResponse().SetError(err), nil
	}
}

func (c *Client) Stop() error {
	c.isStopped = true
	for _, client := range c.clients {
		_ = client.Close()
	}
	return nil
}
