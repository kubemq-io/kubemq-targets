package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"os"
)

type Client struct {
	name   string
	opts   options
	client *storage.Client
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
	b := []byte(c.opts.credentials)

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(b))
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
	case "upload":
		return c.upload(ctx, meta)
	case "download":
		return c.download(ctx, meta)
	case "delete":
		return c.delete(ctx, meta)
	case "list":
		return c.list(ctx, meta)
	case "rename":
		return c.rename(ctx, meta)
	case "copy":
		return c.copy(ctx, meta)
	case "move":
		return c.move(ctx, meta)
	case "create_bucket":
		return c.createBucket(ctx, meta)
	default:
		return nil, errors.New("invalid method type")
	}
}

func (c *Client) createBucket(ctx context.Context, meta metadata) (*types.Response, error) {
	bucket := c.client.Bucket(meta.bucket)

	if err := bucket.Create(ctx, meta.projectID, &storage.BucketAttrs{
		StorageClass: meta.storageClass,
		Location:     meta.location,
	}); err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) upload(ctx context.Context, meta metadata) (*types.Response, error) {
	f, err := os.Open(meta.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	wc := c.client.Bucket(meta.bucket).Object(meta.object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return nil, err
	}
	if err := wc.Close(); err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) download(ctx context.Context, meta metadata) (*types.Response, error) {
	rc, err := c.client.Bucket(meta.bucket).Object(meta.object).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(data),
		nil
}

func (c *Client) delete(ctx context.Context, meta metadata) (*types.Response, error) {
	err := c.client.Bucket(meta.bucket).Object(meta.object).Delete(ctx)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) list(ctx context.Context, meta metadata) (*types.Response, error) {
	it := c.client.Bucket(meta.bucket).Objects(ctx, nil)
	var attrs []*storage.ObjectAttrs
	for {
		attr, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		attrs = append(attrs, attr)
	}
	if len(attrs) == 0 {
		return nil, fmt.Errorf("received 0 objects from list for bucket: %s", meta.bucket)
	}
	b, err := json.Marshal(attrs)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) rename(ctx context.Context, meta metadata) (*types.Response, error) {
	src := c.client.Bucket(meta.bucket).Object(meta.object)
	dst := c.client.Bucket(meta.bucket).Object(meta.renameObject)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return nil, err
	}
	if err := src.Delete(ctx); err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) copy(ctx context.Context, meta metadata) (*types.Response, error) {
	src := c.client.Bucket(meta.bucket).Object(meta.object)
	dst := c.client.Bucket(meta.dstBucket).Object(meta.renameObject)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) move(ctx context.Context, meta metadata) (*types.Response, error) {
	src := c.client.Bucket(meta.bucket).Object(meta.object)
	dst := c.client.Bucket(meta.dstBucket).Object(meta.renameObject)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return nil, err
	}
	if err := src.Delete(ctx); err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) Stop() error {
	return c.client.Close()
}
