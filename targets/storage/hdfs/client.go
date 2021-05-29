package hdfs

import (
	"context"
	"errors"
	hdfs "github.com/colinmarc/hdfs/v2"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
	"github.com/kubemq-hub/kubemq-targets/types"
	"io/ioutil"
)

type Client struct {
	log    *logger.Logger
	opts   options
	client *hdfs.Client
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
	co := setClientOption(c.opts)
	c.client, err = hdfs.NewClient(co)
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
	case "read_file":
		return c.readFile(meta)
	case "write_file":
		return c.writeFile(meta, req.Data)
	case "remove_file":
		return c.removeFile(meta)
	case "rename_file":
		return c.renameFile(meta)
	case "mkdir":
		return c.makeDir(meta)
	case "stat":
		return c.stat(meta)
	default:
		return nil, errors.New("invalid method type")
	}
}

func (c *Client) writeFile(meta metadata, data []byte) (*types.Response, error) {

	writer, err := c.client.Create(meta.filePath)
	if err != nil {
		return nil, err
	}
	_, err = writer.Write(data)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) makeDir(meta metadata) (*types.Response, error) {
	err := c.client.Mkdir(meta.filePath, meta.fileMode)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) removeFile(meta metadata) (*types.Response, error) {
	err := c.client.Remove(meta.filePath)
	if err != nil {
		return nil, err
	}

	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) renameFile(meta metadata) (*types.Response, error) {
	err := c.client.Rename(meta.oldFilePath, meta.filePath)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) stat(meta metadata) (*types.Response, error) {
	file, err := c.client.Stat(meta.filePath)
	if err != nil {
		return nil, err
	}
	b, err := createStatAsByteArray(file)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetData(b).
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) readFile(meta metadata) (*types.Response, error) {
	file, err := c.client.Open(meta.filePath)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetData(bytes).
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) Stop() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

func setClientOption(opts options) hdfs.ClientOptions {
	c := hdfs.ClientOptions{}
	if opts.address != "" {
		c.Addresses = append(c.Addresses, opts.address)
	}
	if opts.user != "" {
		c.User = opts.user
	}
	return c
}
