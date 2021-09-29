package amazonmq

import (
	"context"
	"crypto/tls"
	"github.com/go-stomp/stomp"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
)

type Client struct {
	log  *logger.Logger
	opts options
	conn *stomp.Conn
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

	netConn, err := tls.Dial("tcp", c.opts.host, &tls.Config{})
	if err != nil {
		return err
	}

	c.conn, err = stomp.Connect(netConn, stomp.ConnOpt.Login(c.opts.username, c.opts.password))
	if err != nil {
		return err
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
	err := c.conn.Send(meta.destination, "text/plain", req.Data)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) Stop() error {
	return c.conn.Disconnect()
}
