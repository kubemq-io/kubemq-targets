package etcd

import (
	"context"
	"fmt"
	"google.golang.org/grpc"

	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/pkg/logger"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
	"go.etcd.io/etcd/clientv3"
)

const ()

type Client struct {
	name   string
	opts   options
	client *clientv3.Client
	log    *logger.Logger
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
	clientConfig := clientv3.Config{
		Endpoints:   c.opts.endpoints,
		DialTimeout: c.opts.dialTimout,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}
	c.client, err = clientv3.New(clientConfig)
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

	}
	return nil, nil
}

func (c *Client) Get(ctx context.Context, meta metadata) (*types.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, c.opts.operationTimeout)
	defer cancel()
	resp, err := c.client.Get(ctx, meta.key, clientv3.WithSort(clientv3.SortByVersion, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return types.NewResponse().
			SetMetadataKeyValue("key", meta.key).
			SetMetadataKeyValue("error", "true").
			SetMetadataKeyValue("message", "no data found for this key"), nil
	}

	return types.NewResponse().
		SetData(resp.Kvs[0].Value).
		SetMetadataKeyValue("error", "false").
		SetMetadataKeyValue("key", meta.key), nil
}

func (c *Client) Set(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, c.opts.operationTimeout)
	defer cancel()
	_, err := c.client.Put(ctx, meta.key, string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to set key %s: %s", meta.key, err)
	}
	return types.NewResponse().
			SetMetadataKeyValue("key", meta.key).
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) Delete(ctx context.Context, meta metadata) (*types.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, c.opts.operationTimeout)
	defer cancel()
	_, err := c.client.Delete(ctx, meta.key)
	if err != nil {
		return nil, fmt.Errorf("failed to delete key %s: %s", meta.key, err)
	}
	return types.NewResponse().
			SetMetadataKeyValue("key", meta.key).
			SetMetadataKeyValue("result", "ok"),
		nil
}
