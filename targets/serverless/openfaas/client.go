package openfaas

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"io"
	"io/ioutil"
)

type Client struct {
	name   string
	opts   options
	client *resty.Client
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
