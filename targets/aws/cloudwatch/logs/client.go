package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type Client struct {
	name   string
	opts   options
	client *cloudwatchlogs.CloudWatchLogs
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

	svc := cloudwatchlogs.New(sess)
	c.client = svc

	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	default:
		return nil, fmt.Errorf(getValidMethodTypes())
	}
}

func (c *Client) getLogEvents(ctx context.Context, meta metadata) (*types.Response, error) {
	resp, err := c.client.GetLogEventsWithContext(ctx, &cloudwatchlogs.GetLogEventsInput{
		Limit:         aws.Int64(meta.limit),
		LogGroupName:  aws.String(meta.groupName),
		LogStreamName: aws.String(meta.streamName),
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) createLogEventGroup(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	m := make(map[string]*string)
	if data != nil {
		err := json.Unmarshal(data, &m)
		if err != nil {
			return nil, err
		}
	}
	resp, err := c.client.CreateLogGroupWithContext(ctx, &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(meta.groupName),
		Tags:         m,
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}


func (c *Client) deleteLogEventGroup(ctx context.Context, meta metadata) (*types.Response, error) {
	resp, err := c.client.DeleteLogGroupWithContext(ctx, &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: aws.String(meta.groupName),
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) describeLogGroup(ctx context.Context, meta metadata) (*types.Response, error) {
	resp, err := c.client.DescribeLogGroupsWithContext(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(meta.groupPrefix),
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) listTagsLogGroup(ctx context.Context, meta metadata) (*types.Response, error) {
	resp, err := c.client.ListTagsLogGroupWithContext(ctx, &cloudwatchlogs.ListTagsLogGroupInput{
		LogGroupName: aws.String(meta.groupPrefix),
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}