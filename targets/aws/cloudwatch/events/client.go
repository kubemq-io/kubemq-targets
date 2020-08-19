package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type Client struct {
	name   string
	opts   options
	client *cloudwatchevents.CloudWatchEvents
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

	svc := cloudwatchevents.New(sess)
	c.client = svc

	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "put_targets":
		return c.putTarget(ctx, meta, req.Data)
	case "send_event":
		return c.sendEvent(ctx, meta, req.Data)
	case "list_buses":
		return c.listEventBuses(ctx, meta)
	default:
		return nil, errors.New("invalid method type")
	}
}

func (c *Client) putTarget(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	var m map[string]string
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Targets ,please verify data is map[string]string ,strng:Arn and string:Id")
	}
	var targets []*cloudwatchevents.Target
	for k, v := range m {
		i := cloudwatchevents.Target{
			Arn: aws.String(v),
			Id:  aws.String(k),
		}
		targets = append(targets, &i)
	}
	_, err = c.client.PutTargetsWithContext(ctx, &cloudwatchevents.PutTargetsInput{
		Rule:    aws.String(meta.rule),
		Targets: targets,
	})
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) sendEvent(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	var s []*string
	if data != nil {
		err := json.Unmarshal(data, &s)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Resources ,please verify data is a valid []*string of RESOURCE_ARN ")
		}
	}
	res, err := c.client.PutEventsWithContext(ctx, &cloudwatchevents.PutEventsInput{
		Entries: []*cloudwatchevents.PutEventsRequestEntry{
			{
				Detail:     aws.String(meta.detail),
				DetailType: aws.String(meta.detailType),
				Resources:  s,
				Source:     aws.String(meta.source),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) listEventBuses(ctx context.Context, meta metadata) (*types.Response, error) {
	res, err := c.client.ListEventBusesWithContext(ctx, &cloudwatchevents.ListEventBusesInput{
		Limit: aws.Int64(meta.limit),
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}
