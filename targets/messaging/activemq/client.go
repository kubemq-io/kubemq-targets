package activemq

import (
	"context"
	"fmt"
	"github.com/go-stomp/stomp"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/types"

	"time"
)

const (
	defaultConnectTimeout = 5 * time.Second
)

type Client struct {
	name string
	opts options
	conn *stomp.Conn
}

func New() *Client {
	return &Client{}
}
func (c *Client) Name() string {
	return c.name
}
func (c *Client) Init(ctx context.Context, cfg config.Metadata) error {
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
