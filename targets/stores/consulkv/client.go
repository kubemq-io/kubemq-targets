package consulkv

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
	"github.com/kubemq-hub/kubemq-targets/types"
	"net"
	"strings"
)

// Client is a Client state store
type Client struct {
	log    *logger.Logger
	client *consul.Client
	opts   options
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
	config, err := setConfig(c.opts)
	if err != nil {
		return err
	}
	client, err := consul.NewClient(config)
	if err != nil {
		return err
	}
	c.client = client
	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "get":
		return c.get(ctx, meta)
	case "list":
		return c.list(ctx, meta)
	case "put":
		return c.put(ctx, meta, req.Data)
	case "delete":
		return c.delete(ctx, meta)

	}
	return nil, nil
}

func (c *Client) get(ctx context.Context, meta metadata) (*types.Response, error) {
	if meta.key == "" {
		return nil, errors.New("missing key")
	}
	kv := c.client.KV()
	o := c.createQueryOptions(ctx, meta)

	kvp, qm, err := kv.Get(meta.key, o)
	if err != nil {
		return nil, err
	}
	if kvp == nil {
		return nil, fmt.Errorf("failed to find key %s", meta.key)
	}
	b, err := json.Marshal(kvp)
	if err != nil {
		return nil, err
	}
	qmeta, err := json.Marshal(qm)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetData(b).
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("query_metadata", fmt.Sprintf("%s", qmeta)).
		SetMetadataKeyValue("key", meta.key), nil
}

func (c *Client) put(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	if data == nil {
		return nil, errors.New("missing data")
	}
	kvp := &consul.KVPair{}
	err := json.Unmarshal(data, kvp)
	if err != nil {
		return nil, err
	}
	kv := c.client.KV()
	o := c.createWriteOptions(ctx)
	wm, err := kv.Put(kvp, o)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(wm)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetData(b).
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("key", meta.key), nil
}

func (c *Client) delete(ctx context.Context, meta metadata) (*types.Response, error) {
	if meta.key == "" {
		return nil, errors.New("missing key")
	}
	kv := c.client.KV()
	o := c.createWriteOptions(ctx)

	wm, err := kv.Delete(meta.key, o)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(wm)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetData(b).
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("key", meta.key), nil
}

func (c *Client) list(ctx context.Context, meta metadata) (*types.Response, error) {
	kv := c.client.KV()
	o := c.createQueryOptions(ctx, meta)

	kvp, qm, err := kv.List(meta.prefix, o)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(kvp)
	if err != nil {
		return nil, err
	}
	qmeta, err := json.Marshal(qm)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetData(b).
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("query_metadata", fmt.Sprintf("%s", qmeta)).
		SetMetadataKeyValue("key", meta.key), nil
}

func setConfig(opts options) (*consul.Config, error) {
	c := &consul.Config{}
	if opts.address != "" {
		c.Address = opts.address
	}
	if opts.scheme != "" {
		c.Scheme = opts.scheme
	}
	if opts.datacenter != "" {
		c.Datacenter = opts.datacenter
	}

	if opts.waitTime > 0 {
		c.WaitTime = opts.waitTime
	}
	if opts.token != "" {
		c.Token = opts.token
	}
	if opts.tokenFile != "" {
		c.TokenFile = opts.tokenFile
	}
	if opts.tls {

		tlsClientConfig := &tls.Config{}

		tlsClientConfig.InsecureSkipVerify = opts.insecureSkipVerify

		if opts.address != "" {
			server := opts.address
			hasPort := strings.LastIndex(server, ":") > strings.LastIndex(server, "]")
			if hasPort {
				var err error
				server, _, err = net.SplitHostPort(server)
				if err != nil {
					return nil, err
				}
			}
			tlsClientConfig.ServerName = server
		}

		cert, err := tls.X509KeyPair([]byte(opts.tlscertificatefile), []byte(opts.tlscertificatekey))
		if err != nil {
			return c, fmt.Errorf("consulkv: error parsing client certificate: %v", err)
		}
		tlsClientConfig.Certificates = []tls.Certificate{cert}
	}
	return c, nil
}

func (c *Client) Stop() error {
	return nil
}

func (c *Client) createWriteOptions(ctx context.Context) *consul.WriteOptions {
	o := &consul.WriteOptions{}
	o.WithContext(ctx)
	if c.opts.datacenter != "" {
		o.Datacenter = c.opts.datacenter
	}
	if c.opts.token != "" {
		o.Token = c.opts.token
	}
	return o
}
func (c *Client) createQueryOptions(ctx context.Context, meta metadata) *consul.QueryOptions {
	o := &consul.QueryOptions{}
	o.WithContext(ctx)
	if c.opts.datacenter != "" {
		o.Datacenter = c.opts.datacenter
	}
	o.AllowStale = meta.allowStale
	o.RequireConsistent = meta.requireConsistent
	o.UseCache = meta.useCache
	if meta.maxAge != 0 {
		o.MaxAge = meta.maxAge
	}
	if meta.staleIfError != 0 {
		o.StaleIfError = meta.staleIfError
	}
	o.WaitTime = c.opts.waitTime
	if c.opts.token != "" {
		o.Token = c.opts.token
	}
	o.Near = meta.near
	o.Filter = meta.filter

	return o
}
