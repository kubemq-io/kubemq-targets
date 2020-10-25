package redshift

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type Client struct {
	name   string
	opts   options
	client *redshift.Redshift
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

	svc := redshift.New(sess)
	c.client = svc
	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "create_tags":
		return c.createTags(ctx, meta, req.Data)
	case "delete_tags":
		return c.deleteTags(ctx, meta, req.Data)
	case "list_tags":
		return c.listTags(ctx)
	case "list_snapshots":
		return c.listSnapshots(ctx)
	case "list_snapshots_by_tags_keys":
		return c.listSnapshotsByTagsKeys(ctx, req.Data)
	case "list_snapshots_by_tags_values":
		return c.listSnapshotsTagsValues(ctx, req.Data)
	case "describe_cluster":
		return c.describeCluster(ctx, meta)
	case "list_clusters":
		return c.listClusters(ctx)
	case "list_clusters_by_tags_keys":
		return c.listClustersByTagsKeys(ctx, req.Data)
	case "list_clusters_by_tags_values":
		return c.listClustersByTagsValues(ctx, req.Data)
	}

	return nil, errors.New("invalid method type")
}

func (c *Client) createTags(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	if data == nil {
		return nil, errors.New("missing data,tag list is required")
	}
	tags := make(map[string]string)
	err := json.Unmarshal(data, &tags)
	if err != nil {
		return nil, errors.New("data should be map[string]string,tag key(string),tag value(string)")
	}
	var redshiftTags []*redshift.Tag
	for k, v := range tags {
		t := redshift.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}
		redshiftTags = append(redshiftTags, &t)
	}
	_, err = c.client.CreateTagsWithContext(ctx, &redshift.CreateTagsInput{
		ResourceName: aws.String(meta.resourceARN),
		Tags:         redshiftTags,
	})
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) deleteTags(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	if data == nil {
		return nil, errors.New("missing data , tag list is required")
	}
	var tags []*string
	err := json.Unmarshal(data, &tags)
	if err != nil {
		return nil, errors.New("data should be []*string")
	}
	_, err = c.client.DeleteTagsWithContext(ctx, &redshift.DeleteTagsInput{
		ResourceName: aws.String(meta.resourceARN),
		TagKeys:      tags,
	})
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) listTags(ctx context.Context) (*types.Response, error) {
	m, err := c.client.DescribeTagsWithContext(ctx, nil)
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

func (c *Client) listSnapshots(ctx context.Context) (*types.Response, error) {
	m, err := c.client.DescribeClusterSnapshotsWithContext(ctx, nil)
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

func (c *Client) listSnapshotsByTagsKeys(ctx context.Context, data []byte) (*types.Response, error) {
	if data == nil {
		return nil, errors.New("missing data,tag list keys is required")
	}
	var tags []*string
	err := json.Unmarshal(data, &tags)
	if err != nil {
		return nil, errors.New("data should be []*string")
	}
	m, err := c.client.DescribeClusterSnapshotsWithContext(ctx, &redshift.DescribeClusterSnapshotsInput{
		TagKeys: tags,
	})
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

func (c *Client) listSnapshotsTagsValues(ctx context.Context, data []byte) (*types.Response, error) {
	if data == nil {
		return nil, errors.New("missing data,tag list values is required")
	}
	var tags []*string
	err := json.Unmarshal(data, &tags)
	if err != nil {
		return nil, errors.New("data should be []*string")
	}
	m, err := c.client.DescribeClusterSnapshotsWithContext(ctx, &redshift.DescribeClusterSnapshotsInput{
		TagValues: tags,
	})
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

func (c *Client) describeCluster(ctx context.Context, meta metadata) (*types.Response, error) {
	m, err := c.client.DescribeClustersWithContext(ctx, &redshift.DescribeClustersInput{
		ClusterIdentifier: aws.String(meta.resourceName),
	})
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

func (c *Client) listClusters(ctx context.Context) (*types.Response, error) {
	m, err := c.client.DescribeClustersWithContext(ctx, nil)
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

func (c *Client) listClustersByTagsKeys(ctx context.Context, data []byte) (*types.Response, error) {
	if data == nil {
		return nil, errors.New("missing data,list of tags keys is required")
	}
	var tags []*string
	err := json.Unmarshal(data, &tags)
	if err != nil {
		return nil, errors.New("data should be []*string")
	}
	m, err := c.client.DescribeClustersWithContext(ctx, &redshift.DescribeClustersInput{
		TagKeys: tags,
	})
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

func (c *Client) listClustersByTagsValues(ctx context.Context, data []byte) (*types.Response, error) {
	if data == nil {
		return nil, errors.New("missing data,tag list is required")
	}
	var tags []*string
	err := json.Unmarshal(data, &tags)
	if err != nil {
		return nil, errors.New("data should be []*string")
	}
	m, err := c.client.DescribeClustersWithContext(ctx, &redshift.DescribeClustersInput{
		TagValues: tags,
	})
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

