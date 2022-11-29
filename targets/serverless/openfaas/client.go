package openfaas

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/go-resty/resty/v2"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
)

type Client struct {
	log    *logger.Logger
	opts   options
	client *resty.Client
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
	c.client = resty.New()
	c.client.SetDoNotParseResponse(true)
	c.client.SetBasicAuth(c.opts.username, c.opts.password)
	return nil
}

func readBody(data io.ReadCloser) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	b, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}
	return b, data.Close()
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/%s", c.opts.gateway, meta.topic)
	resp, err := c.client.R().SetContext(ctx).SetBody(req.Data).Post(url)
	if err != nil {
		return nil, err
	}
	body, err := readBody(resp.RawBody())
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetMetadataKeyValue("status", fmt.Sprintf("%d", resp.RawResponse.StatusCode)).
		SetData(body), nil
}

func (c *Client) Stop() error {
	return nil
}
