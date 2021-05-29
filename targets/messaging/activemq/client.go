package activemq

import (
	"context"
	"fmt"
	"github.com/go-stomp/stomp"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
	"github.com/kubemq-hub/kubemq-targets/types"
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

	var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(c.opts.username, c.opts.password),
		stomp.ConnOpt.Host("/"),
	}
	c.conn, err = stomp.Dial("tcp", c.opts.host, options...)
	if err != nil {
		return fmt.Errorf("error connecting to activemq broker, %w", err)
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
	if c.conn != nil {
		return c.conn.Disconnect()
	}
	return nil
}
