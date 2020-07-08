package elastic

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/pkg/logger"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
	"github.com/olivere/elastic/v7"
)

const ()

// Client is a Client state store
type Client struct {
	name    string
	elastic *elastic.Client
	opts    options
	log     *logger.Logger
}

func New() *Client {
	return &Client{}
}
func (c *Client) Name() string {
	return c.name
}
func (c *Client) Init(ctx context.Context, cfg config.Metadata) error {
	c.name = cfg.Name
	c.log = logger.NewLogger(cfg.Name)
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}

	var elasticOpts []elastic.ClientOptionFunc
	elasticOpts = append(elasticOpts,
		elastic.SetURL(c.opts.urls...),
		elastic.SetErrorLog(c.log),
		elastic.SetSniff(c.opts.sniff),
		elastic.SetBasicAuth(c.opts.username, c.opts.password))
	if c.opts.retryBackoffSeconds > 0 {
		retry := elastic.NewBackoffRetrier(elastic.NewSimpleBackoff(c.opts.retryBackoffSeconds * 1000))
		elasticOpts = append(elasticOpts, elastic.SetRetrier(retry))
	}

	c.elastic, err = elastic.NewClient(elasticOpts...)
	if err != nil {
		return err
	}

	return err
}
func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "get":
		return c.Get(ctx, meta)
	case "set":
		return c.Set(ctx, meta, req.Data)
	case "delete":
		return c.Delete(ctx, meta)

	}
	return nil, nil
}
func (c *Client) Get(ctx context.Context, meta metadata) (*types.Response, error) {
	getResp, err := c.elastic.Get().Index(meta.index).Id(meta.id).Do(ctx)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetData(getResp.Source).
		SetMetadataKeyValue("id", meta.id), nil
}

func (c *Client) Set(ctx context.Context, meta metadata, value []byte) (*types.Response, error) {
	setResp, err := c.elastic.Index().Index(meta.index).Id(meta.id).BodyString(string(value)).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to set document id %s: %s", meta.id, err)
	}
	return types.NewResponse().
			SetMetadataKeyValue("id", setResp.Id).
			SetMetadataKeyValue("result", setResp.Result),
		nil
}

func (c *Client) Delete(ctx context.Context, meta metadata) (*types.Response, error) {
	delResp, err := c.elastic.Delete().Index(meta.index).Id(meta.id).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete id '%s',%w", meta.id, err)
	}
	return types.NewResponse().
			SetMetadataKeyValue("id", delResp.Id).
			SetMetadataKeyValue("result", delResp.Result),
		nil
}
