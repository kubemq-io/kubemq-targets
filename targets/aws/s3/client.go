package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type Client struct {
	name       string
	opts       options
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
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
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(c.opts.region),
		Credentials: credentials.NewStaticCredentials(c.opts.awsKey, c.opts.awsSecretKey, c.opts.token),
	})
	if err != nil {
		return err
	}

	svc := s3.New(sess)
	c.client = svc
	c.downloader = s3manager.NewDownloader(sess)
	c.uploader = s3manager.NewUploader(sess)
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
	case "list_bucket_items":
		return c.listBucketItems(ctx, meta)
	case "create_bucket":
		return c.createBucket(ctx, meta)
	case "delete_bucket":
		return c.deleteBucket(ctx, meta)
	case "delete_item_from_bucket":
		return c.deleteItemFromBucket(ctx, meta)
	case "delete_all_items_from_bucket":
		return c.deleteAllItemsFromBucket(ctx, meta)
	case "upload_item":
		return c.uploadItem(ctx, meta, req.Data)
	case "copy_item":
		return c.copyItem(ctx, meta)
	case "get_item":
		return c.downloadItem(ctx, meta)

	default:
		return nil, errors.New("invalid method type")
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
	if meta.waitForCompletion {
		err = c.client.WaitUntilBucketExistsWithContext(ctx, &s3.HeadBucketInput{
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
	_, err := c.client.DeleteBucketWithContext(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(meta.bucketName),
	})
	if err != nil {
		return nil, err
	}
	if meta.waitForCompletion {
		err = c.client.WaitUntilBucketNotExistsWithContext(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(meta.bucketName),
		})
		if err != nil {
			return nil, err
		}
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) deleteItemFromBucket(ctx context.Context, meta metadata) (*types.Response, error) {
	m, err := c.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(meta.bucketName),
		Key:    aws.String(meta.itemName),
	})
	if err != nil {
		return nil, err
	}
	if meta.waitForCompletion {
		err = c.client.WaitUntilObjectNotExistsWithContext(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(meta.bucketName),
			Key:    aws.String(meta.itemName),
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

func (c *Client) deleteAllItemsFromBucket(ctx context.Context, meta metadata) (*types.Response, error) {
	iter := s3manager.NewDeleteListIterator(c.client, &s3.ListObjectsInput{
		Bucket: aws.String(meta.bucketName),
	})
	if err := s3manager.NewBatchDeleteWithClient(c.client).Delete(ctx, iter); err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) uploadItem(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	if c.uploader == nil {
		return nil, fmt.Errorf("uploader client is nil, set uploader to true when creating the client")
	}

	r := bytes.NewReader(data)
	m, err := c.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(meta.bucketName),
		Key:    aws.String(meta.itemName),
		Body:   r,
	})
	if err != nil {
		return nil, err
	}
	if meta.waitForCompletion {
		err = c.client.WaitUntilObjectExistsWithContext(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(meta.bucketName),
			Key:    aws.String(meta.itemName),
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

func (c *Client) copyItem(ctx context.Context, meta metadata) (*types.Response, error) {
	m, err := c.client.CopyObjectWithContext(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(meta.bucketName),
		CopySource: aws.String(meta.copySource),
		Key:        aws.String(meta.itemName),
	})
	if err != nil {
		return nil, err
	}
	if meta.waitForCompletion {
		err = c.client.WaitUntilObjectExistsWithContext(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(meta.bucketName),
			Key:    aws.String(meta.itemName),
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

func (c *Client) downloadItem(ctx context.Context, meta metadata) (*types.Response, error) {
	if c.downloader == nil {
		return nil, fmt.Errorf("downloader client is nil, set downloader to true when creating the client")
	}
	requestInput := s3.GetObjectInput{
		Bucket: aws.String(meta.bucketName),
		Key:    aws.String(meta.itemName),
	}
	buf := aws.NewWriteAtBuffer([]byte{})
	m, err := c.downloader.DownloadWithContext(ctx, buf, &requestInput)
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
func (c *Client) Stop() error {
	return nil
}
