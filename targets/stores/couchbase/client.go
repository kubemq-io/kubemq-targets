package couchbase

import (
	"context"
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
)

const defaultWaitBucketReady = 5 * time.Second

type Client struct {
	name       string
	opts       options
	cluster    *gocb.Cluster
	bucket     *gocb.Bucket
	collection *gocb.Collection
}

func New() *Client {
	return &Client{}
}
func (c *Client) Name() string {
	return c.name
}
func (c *Client) Init(ctx context.Context, cfg config.Metadata) error {
	c.name = cfg.Name
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	clusterOpts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: c.opts.username,
			Password: c.opts.password,
		},
	}
	c.cluster, err = gocb.Connect(c.opts.url, clusterOpts)
	if err != nil {
		return fmt.Errorf("couchbase error: unable to connect to couchbase at %s - %v ", c.opts.url, err)
	}
	c.bucket = c.cluster.Bucket(c.opts.bucket)
	err = c.bucket.WaitUntilReady(defaultWaitBucketReady, nil)
	if err != nil {
		return fmt.Errorf("couchbase error: unable to connect to bucket %s - %v ", c.opts.bucket, err)
	}
	if c.opts.collection == "" {
		c.collection = c.bucket.DefaultCollection()
	} else {
		c.collection = c.bucket.Collection(c.opts.collection)
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
	var timeout time.Duration
	deadline, ok := ctx.Deadline()
	if ok {
		timeout = time.Until(deadline)
	}
	result, err := c.collection.Get(meta.key, &gocb.GetOptions{
		WithExpiry:    false,
		Project:       nil,
		Transcoder:    gocb.NewRawBinaryTranscoder(),
		Timeout:       timeout,
		RetryStrategy: nil,
	})
	if err != nil {
		return nil, err
	}
	var value []byte
	err = result.Content(&value)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
		SetData(value).
		SetMetadataKeyValue("error", "false").
		SetMetadataKeyValue("key", meta.key), nil
}

func (c *Client) Set(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	var timeout time.Duration
	deadline, ok := ctx.Deadline()
	if ok {
		timeout = time.Until(deadline)
	}

	_, err := c.collection.Upsert(meta.key, data, &gocb.UpsertOptions{
		Expiry:        meta.expiry,
		PersistTo:     uint(c.opts.numToPersist),
		ReplicateTo:   uint(c.opts.numToReplicate),
		Transcoder:    gocb.NewRawBinaryTranscoder(),
		Timeout:       timeout,
		RetryStrategy: nil,
	})
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("key", meta.key).
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) Delete(ctx context.Context, meta metadata) (*types.Response, error) {
	var timeout time.Duration
	deadline, ok := ctx.Deadline()
	if ok {
		timeout = time.Until(deadline)
	}
	_, err := c.collection.Remove(meta.key, &gocb.RemoveOptions{
		Cas:             gocb.Cas(meta.cas),
		PersistTo:       uint(c.opts.numToPersist),
		ReplicateTo:     uint(c.opts.numToReplicate),
		DurabilityLevel: 0,
		Timeout:         timeout,
		RetryStrategy:   nil,
	})
	if err != nil {
		return nil, err
	}

	return types.NewResponse().
			SetMetadataKeyValue("key", meta.key).
			SetMetadataKeyValue("result", "ok"),
		nil
}
