package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type Client struct {
	name   string
	opts   options
	client *s3.S3
	log    *logger.Logger
}

func New() *Client {
	return &Client{}

}
func (c *Client) Name() string {
	return c.name
}

func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	c.name = cfg.Name
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(c.opts.region),
		Credentials: credentials.NewStaticCredentials(c.opts.awsKey, c.opts.awsSecretKey, c.opts.token),
	})
	if err != nil {
		return err
	}

	svc := s3.New(sess)
	c.client = svc

	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "list_buckets":
		return c.listBuckets(ctx)
	case "create_bucket":
		return c.createBucket(ctx, meta)

	default:
		return nil, fmt.Errorf(getValidMethodTypes())
	}
}

func (c *Client) listBuckets(ctx context.Context) (*types.Response, error) {
	m, err := c.client.ListBucketsWithContext(ctx, nil)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) listBucketItems(ctx context.Context, meta metadata) (*types.Response, error) {
	m, err := c.client.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{Bucket: aws.String(meta.bucketName)})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) createBucket(ctx context.Context, meta metadata) (*types.Response, error) {
	m, err := c.client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(meta.bucketName),
	})
	if err != nil {
		return nil, err
	}
	if meta.waitForCompletion == true {
		err = c.client.WaitUntilBucketExistsWithContext(ctx,&s3.HeadBucketInput{
			Bucket: aws.String(meta.bucketName),
		})
		if err != nil {
			return nil, err
		}
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}


func (c *Client) deleteBucket(ctx context.Context, meta metadata) (*types.Response, error) {
	m, err := c.client.DeleteBucketWithContext(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(meta.bucketName),
	})
	if err != nil {
		return nil, err
	}
	if meta.waitForCompletion == true {
		err = c.client.WaitUntilBucketNotExistsWithContext(ctx,&s3.HeadBucketInput{
			Bucket: aws.String(meta.bucketName),
		})
		if err != nil {
			return nil, err
		}
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}
