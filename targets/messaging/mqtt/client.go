package mqtt

import (
	"context"
	"fmt"
	"go.uber.org/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
)

const (
	defaultConnectTimeout = 5 * time.Second
)

type Client struct {
	log            *logger.Logger
	opts           options
	client         mqtt.Client
	isConnected    *atomic.Bool
	reconnectCount *atomic.Int32
}

func New() *Client {
	return &Client{
		isConnected:    atomic.NewBool(false),
		reconnectCount: atomic.NewInt32(0),
	}
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
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", c.opts.host))
	opts.SetUsername(c.opts.username)
	opts.SetPassword(c.opts.password)
	opts.SetClientID(c.opts.clientId)
	opts.SetKeepAlive(2)
	opts.SetConnectTimeout(defaultConnectTimeout)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetryInterval(1 * time.Second)
	opts.SetMaxReconnectInterval(24 * time.Hour)
	opts.SetConnectRetry(true)
	opts.SetOnConnectHandler(c.onConnect)
	opts.SetConnectionLostHandler(c.onConnectionLost)
	opts.SetReconnectingHandler(c.onReconnectingHandler)

	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("error connecting to mqtt broker, %w", token.Error())
	}

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
	if c.isConnected.Load() == false {
		c.log.Errorf("publish message to topic %s failed, mqtt client is not connected", meta.topic)
		return nil, fmt.Errorf("mqtt client is not connected")
	}
	c.log.Infof("publish message to topic: %s , with qos %d, payload size: %d, payload: %s", meta.topic, meta.qos, len(req.Data), req.String())
	token := c.client.Publish(meta.topic, byte(meta.qos), false, req.Data)
	token.WaitTimeout(time.Second)
	if token.Error() != nil {
		c.log.Errorf("publish message to topic %s failed, error: %s", meta.topic, token.Error().Error())
		return nil, token.Error()
	}
	c.log.Infof("publish message to topic: %s , with qos %d, paload size: %d", meta.topic, meta.qos, len(req.Data))
	return types.NewResponse().SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) Stop() error {
	if c.client != nil {
		c.log.Info("client stopping")
		c.client.Disconnect(0)
	}
	return nil
}

func (c *Client) onConnectionLost(client mqtt.Client, err error) {
	c.log.Errorf("mqtt client connection lost, error: %s", err.Error())
	c.isConnected.Store(false)

}

func (c *Client) onConnect(client mqtt.Client) {
	c.log.Infof("mqtt client connected")
	c.isConnected.Store(true)
	c.reconnectCount.Store(0)
}

func (c *Client) onReconnectingHandler(client mqtt.Client, opts *mqtt.ClientOptions) {
	c.reconnectCount.Inc()
	c.log.Warnf("mqtt client reconnecting to broker, attempt: %d", c.reconnectCount.Load())

}
