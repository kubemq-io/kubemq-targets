package aerospike

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	aero "github.com/aerospike/aerospike-client-go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type Client struct {
	log    *logger.Logger
	opts   options
	client *aero.Client
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
	c.client, err = aero.NewClient(c.opts.host, c.opts.port)
	if err != nil {
		return fmt.Errorf("error in creating aerospike client: %s", err)
	}
	if c.opts.username != "" {
		err = c.client.CreateUser(&aero.AdminPolicy{
			Timeout: c.opts.timeout,
		}, c.opts.username, c.opts.password, nil)
		if err != nil {
			return err
		}
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
	case "get_batch":
		return c.GetBatch(meta, req.Data)
	case "delete":
		return c.Delete(meta)

	}
	return nil, nil
}

//
func (c *Client) Get(meta metadata) (*types.Response, error) {
	key, err := aero.NewKey(meta.namespace, meta.key, meta.userKey)
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
		SetMetadataKeyValue("user_key", meta.key), nil
}

//
func (c *Client) Put(meta metadata, data []byte) (*types.Response, error) {
	if data == nil {
		return nil, errors.New("missing data PutRequest")
	}
	var kb PutRequest
	err := json.Unmarshal(data, &kb)
	if err != nil {
		return nil, err
	}
	if kb.Namespace == "" && meta.namespace != "" {
		kb.Namespace = meta.namespace
	}
	key, err := aero.NewKey(kb.Namespace, kb.KeyName, kb.UserKey)
	if err != nil {
		return nil, err
	}
	err = c.client.Put(nil, key, kb.BinMap)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("key", kb.KeyName), nil
}

func (c *Client) GetBatch(meta metadata, data []byte) (*types.Response, error) {
	if data == nil {
		return nil, errors.New("missing data GetBatchRequest")
	}
	var kb GetBatchRequest
	err := json.Unmarshal(data, &kb)
	if err != nil {
		return nil, err
	}
	var keys []*aero.Key
	for _, k := range kb.KeyNames {
		if kb.Namespace == "" && meta.namespace != "" {
			kb.Namespace = meta.namespace
		}
		key, err := aero.NewKey(kb.Namespace, *k, meta.userKey)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	rec, err := c.client.BatchGet(nil, keys, kb.BinNames...)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(rec)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetData(b).
		SetMetadataKeyValue("result", "ok"), nil
}

//
func (c *Client) Delete(meta metadata) (*types.Response, error) {
	key, err := aero.NewKey(meta.namespace, meta.key, meta.userKey)
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
