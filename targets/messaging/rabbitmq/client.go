package rabbitmq

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/streadway/amqp"
	"strings"
	"sync"
)

type Client struct {
	sync.Mutex
	log         *logger.Logger
	opts        options
	channel     *amqp.Channel
	conn        *amqp.Connection
	isConnected bool
}

func New() *Client {
	return &Client{
		opts:    options{},
		channel: nil,
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
	if err := c.connect(); err != nil {
		return err
	}
	return nil
}
func (c *Client) getTLSConfig() (*tls.Config, error) {
	tlsCfg := &tls.Config{
		InsecureSkipVerify: c.opts.insecure,
	}
	if c.opts.caCert != "" {
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM([]byte(c.opts.caCert)) {
			return nil, fmt.Errorf("error loading Root CA Cert")
		}
		tlsCfg.RootCAs = caCertPool
		c.log.Infof("TLS CA Cert Loaded for RabbitMQ Connection")
	}
	if c.opts.clientCertificate != "" && c.opts.clientKey != "" {
		cert, err := tls.X509KeyPair([]byte(c.opts.clientCertificate), []byte(c.opts.clientKey))
		if err != nil {
			return nil, fmt.Errorf("error loading tls client key pair, %s", err.Error())
		}
		tlsCfg.Certificates = []tls.Certificate{cert}
		c.log.Infof("TLS Client Key Pair Loaded for RabbitMQ Connection")

	}
	return tlsCfg, nil
}
func (c *Client) connect() error {

	if strings.HasPrefix(c.opts.url, "amqps://") {
		tlsCfg, err := c.getTLSConfig()
		if err != nil {
			return err
		}
		c.conn, err = amqp.DialTLS(c.opts.url, tlsCfg)
		if err != nil {
			return fmt.Errorf("error dialing rabbitmq, %w", err)
		}

	} else {
		var err error
		c.conn, err = amqp.Dial(c.opts.url)
		if err != nil {
			return fmt.Errorf("error dialing rabbitmq, %w", err)
		}
	}
	var err error
	c.channel, err = c.conn.Channel()
	if err != nil {
		_ = c.conn.Close()
		return fmt.Errorf("error getting rabbitmq channel, %w", err)
	}
	c.isConnected = true
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
	return c.Publish(ctx, meta, req.Data)
}

func (c *Client) Publish(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	c.Lock()
	defer c.Unlock()
	if !c.isConnected {
		if err := c.connect(); err != nil {
			return nil, err
		}
	}
	msg := meta.amqpMessage(data)
	err := c.channel.Publish(meta.exchange, meta.queue, meta.mandatory, meta.immediate, msg)
	if err != nil {
		c.isConnected = false
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) Stop() error {
	if c.channel != nil {
		return c.channel.Close()
	}
	return nil
}
