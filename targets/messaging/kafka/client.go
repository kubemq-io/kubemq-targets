package kafka

import (
	"context"
	"crypto/tls"
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
	kc.Producer.RequiredAcks = kafka.WaitForAll
	kc.Producer.Retry.Max = 5
	kc.Producer.Return.Successes = true
	if c.opts.saslUsername != "" {
		kc.Net.SASL.Enable = true
		kc.Net.SASL.User = c.opts.saslUsername
		kc.Net.SASL.Password = c.opts.saslPassword

		kc.Net.TLS.Enable = true
		kc.Net.TLS.Config = &tls.Config{
			ClientAuth: 0,
		}
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
