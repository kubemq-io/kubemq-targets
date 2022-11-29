package rethinkdb

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
	rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type Client struct {
	log     *logger.Logger
	opts    options
	session *rethink.Session
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
	opts, err := c.setUpConnectOpts(c.opts)
	if err != nil {
		return err
	}
	c.session, err = rethink.Connect(opts)
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
		return c.get(ctx, meta)
	case "update":
		return c.update(ctx, meta, req.Data)
	case "insert":
		return c.insert(ctx, meta, req.Data)
	case "delete":
		return c.delete(ctx, meta)

	}
	return nil, nil
}

func (c *Client) get(ctx context.Context, meta metadata) (*types.Response, error) {
	cursor, err := rethink.DB(meta.dbName).Table(meta.table).Get(meta.key).Run(c.session, rethink.RunOpts{Context: ctx})
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var resp []interface{}
	var item interface{}
	for cursor.Next(&item) {
		if item != nil {
			resp = append(resp, item)
		}
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("failed to find value for request on table %s", meta.table)
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetData(b).
		SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) insert(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	if data == nil {
		return nil, errors.New("missing data")
	}
	var args map[string]interface{}
	err := json.Unmarshal(data, &args)
	if err != nil {
		return nil, err
	}
	_, err = rethink.DB(meta.dbName).Table(meta.table).Insert(args).Run(c.session, rethink.RunOpts{Context: ctx})
	if err != nil {
		return nil, err
	}

	return types.NewResponse().
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("table", meta.table), nil
}

func (c *Client) update(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	if data == nil {
		return nil, errors.New("missing data")
	}
	var args map[string]interface{}
	err := json.Unmarshal(data, &args)
	if err != nil {
		return nil, err
	}
	_, err = rethink.DB(meta.dbName).Table(meta.table).Get(meta.key).Update(args).Run(c.session, rethink.RunOpts{Context: ctx})
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("key", meta.key), nil
}

func (c *Client) delete(ctx context.Context, meta metadata) (*types.Response, error) {
	_, err := rethink.DB(meta.dbName).Table(meta.table).Get(meta.key).Delete().Run(c.session, rethink.RunOpts{Context: ctx})
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("key", meta.key), nil
}

func (c *Client) setUpConnectOpts(opts options) (rethink.ConnectOpts, error) {
	co := rethink.ConnectOpts{
		Address: opts.host,
	}
	if opts.username != "" {
		co.Username = opts.username
	}
	if opts.password != "" {
		co.Password = opts.password
	}
	if opts.authKey != "" {
		co.AuthKey = opts.authKey
	}
	if opts.timeout != 0 {
		co.Timeout = opts.timeout
	}
	if opts.keepAlivePeriod != 0 {
		co.KeepAlivePeriod = opts.keepAlivePeriod
	}
	if opts.ssl {

		cert, err := tls.X509KeyPair([]byte(opts.certFile), []byte(opts.certKey))
		if err != nil {
			return co, fmt.Errorf("rethinkdb: error parsing session certificate: %v", err)
		}
		co.TLSConfig.Certificates = append(co.TLSConfig.Certificates, cert)
	}
	if opts.handShakeVersion != 0 {
		co.HandshakeVersion = rethink.HandshakeVersion(opts.handShakeVersion)
	}
	if opts.numberOfRetries != 0 {
		co.NumRetries = opts.numberOfRetries
	}
	if opts.initialCap != 0 {
		co.InitialCap = opts.initialCap
	}
	if opts.maxOpen != 0 {
		co.MaxOpen = opts.maxOpen
	}
	return co, nil
}

func (c *Client) Stop() error {
	return c.session.Close(rethink.CloseOpts{})
}
