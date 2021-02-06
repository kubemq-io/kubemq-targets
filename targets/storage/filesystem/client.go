package filesystem

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Client struct {
	opts    options
	absPath string
}

func New() *Client {
	return &Client{}
}
func (c *Client) Connector() *common.Connector {
	return Connector()
}
func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	if _, err := os.Stat(c.opts.basePath); os.IsNotExist(err) {
		return fmt.Errorf("base path is not exist")
	}
	c.absPath, _ = filepath.Abs(c.opts.basePath)
	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "save":
		return c.Save(ctx, meta, req.Data)
	case "load":
		return c.Load(ctx, meta)
	case "delete":
		return c.Delete(ctx, meta)
	case "list":
		return c.List(ctx, meta)
	}
	return nil, nil
}

func (c *Client) Save(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	if _, err := os.Stat(filepath.Join(c.absPath, meta.path)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Join(c.absPath, meta.path), 0600)
		if err != nil {
			return types.NewResponse().SetError(err), nil
		}

	}
	fullPath := filepath.Join(c.absPath, meta.path, meta.filename)
	err := ioutil.WriteFile(fullPath, data, 0600)
	if err != nil {
		return types.NewResponse().SetError(err), nil
	}

	return types.NewResponse().SetMetadataKeyValue("result", "ok"), nil
}
func (c *Client) Delete(ctx context.Context, meta metadata) (*types.Response, error) {
	fullPath := filepath.Join(c.absPath, meta.path, meta.filename)
	err := os.Remove(fullPath)
	if err != nil {
		return types.NewResponse().SetError(err), nil
	}
	return types.NewResponse().SetMetadataKeyValue("result", "ok"), nil
}

func (c *Client) Load(ctx context.Context, meta metadata) (*types.Response, error) {
	fullPath := filepath.Join(c.absPath, meta.path, meta.filename)
	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return types.NewResponse().SetError(err), nil
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(data),
		nil
}
func (c *Client) List(ctx context.Context, meta metadata) (*types.Response, error) {
	var list FileInfoList
	err := filepath.Walk(c.absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		list = append(list, newFromOSFileInfo(info, path))
		return nil
	})
	if err != nil {
		return types.NewResponse().SetError(err), nil
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(list.Marshal()),
		nil
}
func (c *Client) Stop() error {
	return nil
}
