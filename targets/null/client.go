package null

import (
	"context"
	"github.com/kubemq-hub/builder/common"

	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"

	"time"
)

type Client struct {
	name          string
	Delay         time.Duration
	DoError       error
	ResponseError error
}

func (c *Client) Connector() *common.Connector {
	return Connector()
}
func (c *Client) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	select {
	case <-time.After(c.Delay):
		if c.DoError != nil {
			return nil, c.DoError
		}
		if c.ResponseError != nil {
			return nil, c.ResponseError
		}

		return types.NewResponse().SetData(request.Data), nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	c.name = cfg.Name
	return nil
}
