package null

import (
	"context"

	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/types"

	"time"
)

type Client struct {
	name          string
	Delay         time.Duration
	DoError       error
	ResponseError error
}

func (c *Client) Name() string {
	return c.name
}
func (c *Client) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	select {
	case <-time.After(c.Delay):
		if c.DoError != nil {
			return nil, c.DoError
		}
		if c.ResponseError != nil {
			return types.NewResponse().
					SetMetadataKeyValue("error", c.ResponseError.Error()),
				nil
		}

		return types.NewResponse().SetData(request.Data), nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Client) Init(ctx context.Context, cfg config.Metadata) error {
	c.name = cfg.Name
	return nil
}
