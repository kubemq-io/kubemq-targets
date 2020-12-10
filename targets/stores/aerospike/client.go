package aerospike

import (
	"context"
	"encoding/json"
	"fmt"
	aero "github.com/aerospike/aerospike-client-go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type Client struct {
	name   string
	opts   options
	client *aero.Client
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
	c.client, err = aero.NewClient(c.opts.host, c.opts.port)
	if err != nil {
		return fmt.Errorf("error in creating aerospike client: %s", err)
	}
	err = c.client.CreateUser(&aero.AdminPolicy{
		Timeout: c.opts.timeout,
	}, c.opts.username, c.opts.password, nil)
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
		return c.Get(meta)
	case "set":
		return c.Put(meta, req.Data)
	case "delete":
		return c.Delete(meta)
		
	}
	return nil, nil
}

//
func (c *Client) Get(meta metadata) (*types.Response, error) {
	key, err := aero.NewKey(meta.namespace, meta.key, nil)
	if err != nil {
		return nil, err
	}
	rec, err := c.client.Get(nil, key)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, fmt.Errorf("no data found for key %s", key.String())
	}
	b, err := json.Marshal(rec)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetData(b).
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("key", meta.key), nil
}

//
func (c *Client) Put(meta metadata, data []byte) (*types.Response, error) {
	key, err := aero.NewKey(meta.namespace, meta.key, data)
	if err != nil {
		return nil, err
	}
	err = c.client.Put(nil, key, nil)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("key", meta.key), nil
}

//
func (c *Client) Delete(meta metadata) (*types.Response, error) {
	key, err := aero.NewKey(meta.namespace, meta.key, nil)
	if err != nil {
		return nil, err
	}
	del, err := c.client.Delete(nil, key)
	if err != nil {
		return nil, err
	}
	if !del {
		return nil, fmt.Errorf("failed to delete %s", key.String())
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("key", meta.key), nil
}

func (c *Client) Stop() error {
	c.client.Close()
	return nil
}
