package nats

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/nats-io/nats.go"
	"time"
)

type Client struct {
	log    *logger.Logger
	opts   options
	client *nats.Conn
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
	o := setOptions(c.opts.certFile, c.opts.certKey, c.opts.username, c.opts.password, c.opts.token, c.opts.tls, c.opts.timeout)
	c.client, err = nats.Connect(c.opts.url, o)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	err = c.client.Publish(meta.subject, req.Data)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) Stop() error {
	if c.client != nil {
		c.client.Close()
	}
	return nil
}

func setOptions(sslcertificatefile string, sslcertificatekey string, username string, password string, token string, useTls bool, timeout int) nats.Option {
	return func(o *nats.Options) error {
		if useTls {
			if sslcertificatefile != "" && sslcertificatekey != "" {
				cert, err := tls.X509KeyPair([]byte(sslcertificatefile), []byte(sslcertificatekey))
				if err != nil {
					return fmt.Errorf("nats: error parsing client certificate: %v", err)
				}
				cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
				if err != nil {
					return fmt.Errorf("nats: error parsing client certificate: %v", err)
				}
				if o.TLSConfig == nil {
					o.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
				}
				o.TLSConfig.Certificates = []tls.Certificate{cert}
				o.Secure = true
			} else {
				return errors.New("when using tls make sure to pass file and key")
			}
		}
		if username != "" {
			o.User = username
		}
		if password != "" {
			o.Password = password
		}
		if token != "" {
			o.Token = token
		}
		if timeout != 0 {
			o.Timeout = time.Duration(timeout) * time.Second
		}

		return nil
	}
}
