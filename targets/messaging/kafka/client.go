package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strconv"

	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"

	kafka "github.com/Shopify/sarama"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/types"
)

type Client struct {
	log      *logger.Logger
	producer kafka.SyncProducer
	opts     options
	config   *kafka.Config
}

func New() *Client {
	return &Client{}
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

	kc := kafka.NewConfig()
	kc.Version = kafka.V2_0_0_0
	isSSL, isSASL := c.opts.parseSecurityProtocol()

	if isSASL {
		kc.Net.SASL.Enable = isSASL
		kc.Net.SASL.User = c.opts.saslUsername
		kc.Net.SASL.Password = c.opts.saslPassword
		kc.Net.SASL.Mechanism = c.opts.parseASLMechanism()
	}
	if isSSL {
		kc.Net.TLS.Enable = true
		tlsCfg := &tls.Config{
			InsecureSkipVerify: c.opts.insecure,
		}
		if c.opts.cacert != "" {
			caCertPool := x509.NewCertPool()
			if !caCertPool.AppendCertsFromPEM([]byte(c.opts.cacert)) {
				return fmt.Errorf("error loading Root CA Cert")
			}
			tlsCfg.RootCAs = caCertPool
			c.log.Infof("TLS CA Cert Loaded for Kafka Connection")
		}
		if c.opts.clientCert != "" && c.opts.clientKey != "" {
			cert, err := tls.X509KeyPair([]byte(c.opts.clientCert), []byte(c.opts.clientKey))
			if err != nil {
				return fmt.Errorf("error loading tls client key pair, %s", err.Error())
			}
			tlsCfg.Certificates = []tls.Certificate{cert}
			c.log.Infof("TLS Client Key Pair Loaded for Kafka Connection")
		}
		kc.Net.TLS.Config = tlsCfg
	}
	c.config = kc
	c.producer, err = kafka.NewSyncProducer(c.opts.brokers, kc)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	m, err := parseMetadata(request.Metadata, c.opts)
	if err != nil {
		return nil, err
	}

	partition, offset, err := c.producer.SendMessage(&kafka.ProducerMessage{
		Headers: m.Headers,
		Key:     kafka.ByteEncoder(m.Key),
		Value:   kafka.ByteEncoder(request.Data),
		Topic:   c.opts.topic,
	})
	if err != nil {
		return nil, err
	}
	r := types.NewResponse().
		SetMetadataKeyValue("partition", strconv.FormatInt(int64(partition), 10)).
		SetMetadataKeyValue("offset", strconv.FormatInt(offset, 10))
	return r, nil
}

func (c *Client) Connector() *common.Connector {
	return Connector()
}

func (c *Client) Stop() error {
	if c.producer != nil {
		c.config.MetricRegistry.UnregisterAll()
		return c.producer.Close()
	}
	return nil
}
