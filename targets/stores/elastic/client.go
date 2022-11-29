package elastic

import (
	"context"
	"fmt"

	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
	"github.com/olivere/elastic/v7"
)

type Client struct {
	log     *logger.Logger
	elastic *elastic.Client
	opts    options
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

	var elasticOpts []elastic.ClientOptionFunc
	elasticOpts = append(elasticOpts,
		elastic.SetURL(c.opts.urls...),
		elastic.SetSniff(c.opts.sniff),
		elastic.SetBasicAuth(c.opts.username, c.opts.password))

	c.elastic, err = elastic.NewClient(elasticOpts...)
	if err != nil {
		return err
	}
	_, _, err = c.elastic.Ping(c.opts.urls[0]).Do(ctx)
	if err != nil {
		return err
	}

	return nil
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
	case "index.exists":
		return c.IndexExists(ctx, meta)
	case "index.create":
		return c.IndexCreate(ctx, meta, req.Data)
	case "index.delete":
		return c.IndexDelete(ctx, meta)
	default:
		return nil, fmt.Errorf("invalid method")
	}
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

func (c *Client) IndexExists(ctx context.Context, meta metadata) (*types.Response, error) {
	exists, err := c.elastic.IndexExists(meta.index).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute index exist '%s',%w", meta.index, err)
	}
	return types.NewResponse().
			SetMetadataKeyValue("exists", fmt.Sprintf("%t", exists)),
		nil
}

func (c *Client) IndexCreate(ctx context.Context, meta metadata, value []byte) (*types.Response, error) {
	result, err := c.elastic.CreateIndex(meta.index).BodyString(string(value)).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create index'%s',%w", meta.index, err)
	}
	return types.NewResponse().
			SetMetadataKeyValue("acknowledged", fmt.Sprintf("%t", result.Acknowledged)).
			SetMetadataKeyValue("shards_acknowledged", fmt.Sprintf("%t", result.ShardsAcknowledged)).
			SetMetadataKeyValue("index", result.Index),
		nil
}

func (c *Client) IndexDelete(ctx context.Context, meta metadata) (*types.Response, error) {
	result, err := c.elastic.DeleteIndex(meta.index).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete index'%s',%w", meta.index, err)
	}
	return types.NewResponse().
			SetMetadataKeyValue("acknowledged", fmt.Sprintf("%t", result.Acknowledged)),
		nil
}

func (c *Client) Stop() error {
	return nil
}
