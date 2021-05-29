package null

import (
	"context"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"

	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"

	"time"
)

type Client struct {
	log           *logger.Logger
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

func (c *Client) Init(ctx context.Context, cfg config.Spec, log *logger.Logger) error {
	c.log = log
	if c.log == nil {
		c.log = logger.NewLogger(cfg.Kind)
	}

	return nil
}

func (c *Client) Stop() error {
	return nil
}
