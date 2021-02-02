package echo

import (
	"context"
	"github.com/kubemq-hub/builder/connector/common"
	"os"

	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type Client struct {
	host string
}

func New() *Client {
	return &Client{}
}

func (c *Client) Connector() *common.Connector {
	return Connector()
}
func (c *Client) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	return types.NewResponse().
		SetMetadata(request.Metadata).
		SetMetadataKeyValue("host", c.host).
		SetData(request.Data), nil

}

func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	var err error
	c.host, err = os.Hostname()
	if err != nil {
		c.host = "unknown"
	}
	return nil
}

func (c *Client) Stop() error {
	return nil
}
