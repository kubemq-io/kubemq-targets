package activemq

import (
	"context"
	"fmt"
	"github.com/go-stomp/stomp"
	"github.com/kubemq-hub/builder/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type Client struct {
	name string
	opts options
	conn *stomp.Conn
}

func New() *Client {
	return &Client{}
}
func (c *Client) Connector() *common.Connector {
	return Connector()
}
func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	c.name = cfg.Name
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
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	err = c.conn.Send(meta.destination, "text/plain", req.Data)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().SetMetadataKeyValue("result", "ok"), nil
}
