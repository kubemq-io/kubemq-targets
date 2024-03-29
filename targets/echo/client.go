package echo

import (
	"context"
	"fmt"
	"os"

	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"

	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/types"
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
	m, _ := parseMetadata(request.Metadata)
	if m.isError {
		return types.NewResponse().
			SetError(fmt.Errorf("echo error")).
			SetMetadata(request.Metadata).
			SetMetadataKeyValue("host", c.host).
			SetData(request.Data), nil
	}
	return types.NewResponse().
		SetMetadata(request.Metadata).
		SetMetadataKeyValue("host", c.host).
		SetData(request.Data), nil
}

func (c *Client) Init(ctx context.Context, cfg config.Spec, log *logger.Logger) error {
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
