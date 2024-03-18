package amqp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"sync"

	"github.com/Azure/go-amqp"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
)

type Client struct {
	sync.Mutex
	log         *logger.Logger
	opts        options
	session     *amqp.Session
	conn        *amqp.Conn
	senders     map[string]*amqp.Sender
	isConnected bool
}

func New() *Client {
	return &Client{
		opts:    options{},
		session: nil,
		conn:    nil,
		senders: map[string]*amqp.Sender{},
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
	if err := c.connect(ctx); err != nil {
		return err
	}
	return nil
}

func (c *Client) getTLSConfig() (*tls.Config, error) {
	tlsCfg := &tls.Config{
		InsecureSkipVerify: c.opts.insecure,
	}
	if c.opts.insecure {
		c.log.Infof("AMQP connection is configured to skip certificate verification")
	}

	if c.opts.caCert != "" {
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM([]byte(c.opts.caCert)) {
			return nil, fmt.Errorf("error loading Root CA Cert")
		}
		tlsCfg.RootCAs = caCertPool
		c.log.Infof("TLS CA Cert Loaded for RabbitMQ Connection")
	}
	return tlsCfg, nil
}

func (c *Client) connect(ctx context.Context) error {
	connOptions, err := c.opts.getConnOptions()
	c.conn, err = amqp.Dial(ctx, c.opts.url, connOptions)
	if err != nil {
		return fmt.Errorf("error dialing active, %w", err)
	}
	sessionOpt := &amqp.SessionOptions{}
	c.session, err = c.conn.NewSession(ctx, sessionOpt)
	if err != nil {
		_ = c.conn.Close()
		return fmt.Errorf("error getting active session, %w", err)
	}
	c.isConnected = true
	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	return c.Publish(ctx, meta, req.Data)
}

func (c *Client) Publish(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	c.Lock()
	defer c.Unlock()
	if !c.isConnected {
		if err := c.connect(ctx); err != nil {
			return nil, err
		}
	}
	queueSender, ok := c.senders[meta.address]
	if !ok {
		newSender, err := c.session.NewSender(ctx, meta.address, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating new sender, %w", err)
		}
		c.senders[meta.address] = newSender
		queueSender = newSender
	}
	err := queueSender.Send(ctx, meta.amqpMessage(data), nil)
	if err != nil {
		return &types.Response{
			Metadata: nil,
			Data:     nil,
			IsError:  true,
			Error:    err.Error(),
		}, nil
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) Stop() error {
	if c.session != nil {
		return c.conn.Close()
	}
	return nil
}
