package sns

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type Client struct {
	name   string
	opts   options
	client *sns.SNS
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

	svc := sns.New(sess)
	c.client = svc

	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "list_topics":
		return c.listTopics(ctx)

	default:
		return nil, fmt.Errorf(getValidMethodTypes())
	}
}

func (c *Client) listTopics(ctx context.Context) (*types.Response, error) {
	l, err := c.client.ListTopicsWithContext(ctx, nil)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) listSubscriptions(ctx context.Context) (*types.Response, error) {
	l, err := c.client.ListSubscriptionsWithContext(ctx, nil)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) listSubscriptionsByTopic(ctx context.Context, meta metadata) (*types.Response, error) {
	t, err := c.client.ListSubscriptionsByTopicWithContext(ctx, &sns.ListSubscriptionsByTopicInput{
		TopicArn: aws.String(meta.topic),
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) createTopic(ctx context.Context, meta metadata) (*types.Response, error) {
	r, err := c.client.CreateTopicWithContext(ctx, &sns.CreateTopicInput{
		Name: aws.String(meta.topic),
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) subscribeToTopic(ctx context.Context, meta metadata) (*types.Response, error) {
	r, err := c.client.SubscribeWithContext(ctx, &sns.SubscribeInput{
		TopicArn:              aws.String(meta.topic),
		Endpoint:              aws.String(meta.endPoint),
		Protocol:              aws.String(meta.protocol),
		ReturnSubscriptionArn: aws.Bool(meta.returnSubscription),
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) sendingMessageToTopic(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	m := make(map[string]*sns.MessageAttributeValue)
	if data != nil {
		err := json.Unmarshal(data, &m)
		if err != nil {
			return nil, err
		}
	}
	r, err := c.client.PublishWithContext(ctx, &sns.PublishInput{
		TopicArn:          aws.String(meta.topic),
		Message:           aws.String(meta.message),
		PhoneNumber:       aws.String(meta.phoneNumber),
		Subject:           aws.String(meta.subject),
		TargetArn:         aws.String(meta.targetArn),
		MessageAttributes: m,
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) deleteTopic(ctx context.Context, meta metadata) (*types.Response, error) {
	r, err := c.client.DeleteTopicWithContext(ctx, &sns.DeleteTopicInput{
		TopicArn: aws.String(meta.topic),
	})
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}
