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
	case "create_log_event_stream":
		return c.createLogEventStream(ctx, meta)
	case "describe_log_event_stream":
		return c.describeLogEventStream(ctx, meta)
	case "delete_log_event_stream":
		return c.deleteLogEventStream(ctx, meta)
	case "get_log_event":
		return c.getLogEvents(ctx, meta)
	case "create_log_group":
		return c.createLogEventGroup(ctx, meta, req.Data)
	case "describe_log_group":
		return c.describeLogGroup(ctx, meta)
	case "delete_log_group":
		return c.deleteLogEventGroup(ctx, meta)
	case "list_tags_group":
		return c.listTagsLogGroup(ctx, meta)
	case "describe_resources_policy":
		return c.describeResourcePolicies(ctx, meta)
	case "delete_resources_policy":
		return c.deleteResourcePolicies(ctx, meta)
	case "put_resources_policy":
		return c.putResourcePolicies(ctx, meta)
	default:
		return nil, fmt.Errorf(getValidMethodTypes())
	}
}

func (c *Client) createLogEventStream(ctx context.Context, meta metadata) (*types.Response, error) {
	resp, err := c.client.CreateLogStreamWithContext(ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(meta.logGroupName),
		LogStreamName: aws.String(meta.logStreamName),
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

func (c *Client) describeLogEventStream(ctx context.Context, meta metadata) (*types.Response, error) {
	resp, err := c.client.DescribeLogStreamsWithContext(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(meta.logGroupName),
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

func (c *Client) deleteLogEventStream(ctx context.Context, meta metadata) (*types.Response, error) {
	resp, err := c.client.DeleteLogStreamWithContext(ctx, &cloudwatchlogs.DeleteLogStreamInput{
		LogGroupName:  aws.String(meta.logGroupName),
		LogStreamName: aws.String(meta.logStreamName),
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

func (c *Client) getLogEvents(ctx context.Context, meta metadata) (*types.Response, error) {
	resp, err := c.client.GetLogEventsWithContext(ctx, &cloudwatchlogs.GetLogEventsInput{
		Limit:         aws.Int64(meta.limit),
		LogGroupName:  aws.String(meta.logGroupName),
		LogStreamName: aws.String(meta.logStreamName),
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
	resp := &cloudwatchlogs.CreateLogGroupOutput{}
	var err error
	if data != nil {
		err := json.Unmarshal(data, &m)
		if err != nil {
			return nil, err
		}
		resp, err = c.client.CreateLogGroupWithContext(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(meta.logGroupName),
			Tags:         m,
		})
		if err != nil {
			return nil, err
		}
	} else {
		resp, err = c.client.CreateLogGroupWithContext(ctx, &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(meta.logGroupName),
		})
		if err != nil {
			return nil, err
		}
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
		LogGroupName: aws.String(meta.logGroupName),
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
		LogGroupNamePrefix: aws.String(meta.logGroupPrefix),
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
		LogGroupName: aws.String(meta.logGroupPrefix),
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

func (c *Client) describeResourcePolicies(ctx context.Context, meta metadata) (*types.Response, error) {
	resp, err := c.client.DescribeResourcePoliciesWithContext(ctx, &cloudwatchlogs.DescribeResourcePoliciesInput{
		Limit: aws.Int64(meta.limit),
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

func (c *Client) deleteResourcePolicies(ctx context.Context, meta metadata) (*types.Response, error) {
	resp, err := c.client.DeleteResourcePolicyWithContext(ctx, &cloudwatchlogs.DeleteResourcePolicyInput{
		PolicyName: aws.String(meta.policyName),
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

func (c *Client) putResourcePolicies(ctx context.Context, meta metadata) (*types.Response, error) {
	resp, err := c.client.PutResourcePolicyWithContext(ctx, &cloudwatchlogs.PutResourcePolicyInput{
		PolicyName:     aws.String(meta.policyName),
		PolicyDocument: aws.String(meta.policyDocument),
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
