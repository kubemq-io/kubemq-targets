package hazelcast

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/hazelcast/hazelcast-go-client"
	hazelconfig "github.com/hazelcast/hazelcast-go-client/config"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
)

// Client is a Client state store
type Client struct {
	log    *logger.Logger
	client hazelcast.Client
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
	client, err := hazelcast.NewClientWithConfig(config)
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
		return c.get(meta)
	case "get_list":
		return c.getList(meta)
	case "set":
		return c.set(meta, req.Data)
	case "delete":
		return c.delete(meta)

	}
	return nil, nil
}

func (c *Client) get(meta metadata) (*types.Response, error) {
	Map, err := c.client.GetMap(meta.mapName)
	if err != nil {
		return nil, err
	}
	v, err := Map.Get(meta.key)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, fmt.Errorf("could not fine key for value %s", meta.key)
	}
	valueInterface, ok := v.([]byte)
	if !ok {
		return nil, fmt.Errorf("failed to cast interface for key %s to byte array", meta.key)
	}
	return types.NewResponse().
		SetData(valueInterface).
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("key", meta.key), nil
}

func (c *Client) set(meta metadata, value []byte) (*types.Response, error) {
	Map, err := c.client.GetMap(meta.mapName)
	if err != nil {
		return nil, err
	}
	err = Map.Set(meta.key, value)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("key", meta.key), nil
}

func (c *Client) getList(meta metadata) (*types.Response, error) {
	list, err := c.client.GetList(meta.listName)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetData(b).
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("key", meta.key), nil
}

func (c *Client) delete(meta metadata) (*types.Response, error) {
	Map, err := c.client.GetMap(meta.mapName)
	if err != nil {
		return nil, err
	}
	err = Map.Delete(meta.key)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetMetadataKeyValue("result", "ok").
		SetMetadataKeyValue("key", meta.key), nil
}

func setConfig(opts options) (*hazelconfig.Config, error) {
	c := hazelcast.NewConfig()

	if opts.username != "" {
		c.GroupConfig().SetName(opts.username)
	}
	if opts.password != "" {
		c.GroupConfig().SetPassword(opts.password)
	}
	if opts.address != "" {
		c.NetworkConfig().AddAddress(opts.address)
	}
	if opts.connectionAttemptLimit > 0 {
		c.NetworkConfig().SetConnectionAttemptLimit(opts.connectionAttemptLimit)
	}
	if opts.connectionAttemptPeriod > 0 {
		c.NetworkConfig().SetConnectionAttemptPeriod(opts.connectionAttemptPeriod)
	}
	if opts.connectionTimeout > 0 {
		c.NetworkConfig().SetConnectionTimeout(opts.connectionTimeout)
	}
	if opts.ssl {

		cert, err := tls.X509KeyPair([]byte(opts.sslcertificatefile), []byte(opts.sslcertificatekey))
		if err != nil {
			return c, fmt.Errorf("hazelcast: error parsing client certificate: %v", err)
		}
		sslConfig := c.NetworkConfig().SSLConfig()
		sslConfig.SetEnabled(true)
		sslConfig.Certificates = append(sslConfig.Certificates, cert)
		if opts.serverName != "" {
			sslConfig.ServerName = opts.serverName
		}
	}
	return c, nil
}

func (c *Client) Stop() error {
	if c.client != nil {
		c.client.Shutdown()
	}
	return nil
}
