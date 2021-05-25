package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"

	"strconv"
)

type Client struct {
	name   string
	opts   options
	client *sqs.SQS
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
		Credentials: credentials.NewStaticCredentials(c.opts.sqsKey, c.opts.sqsSecretKey, c.opts.token),
	})
	if err != nil {
		return err
	}

	svc := sqs.New(sess)
	c.client = svc

	return nil
}

func (c *Client) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	eventMetadata, ok := c.opts.defaultMetadata()
	if !ok {
		var err error
		eventMetadata, err = parseMetadata(request.Metadata, c.opts)
		if err != nil {
			return nil, err
		}
	}
	m := &sqs.SendMessageInput{}
	c.setMessageMeta(m, eventMetadata)
	m.SetMessageBody(fmt.Sprintf("%s", request.Data))
	tries := 0
	for tries <= c.opts.retries {
		r, err := c.client.SendMessageWithContext(ctx, m)
		if err == nil {
			return types.NewResponse().
					SetMetadataKeyValue("event_id", *r.MessageId),
				nil
		}
		if tries >= c.opts.retries {
			return nil, err
		}
		tries++
	}
	return nil, fmt.Errorf("retries must be a zero or greater")
}

func (c *Client) setMessageMeta(m *sqs.SendMessageInput, eventMetadata metadata) *sqs.SendMessageInput {
	m.SetQueueUrl(eventMetadata.queueURL)
	m.SetDelaySeconds(int64(eventMetadata.delay))
	if len(eventMetadata.tags) > 0 {
		m.SetMessageAttributes(eventMetadata.tags)
	}
	return m
}

func (c *Client) SetQueueAttributes(ctx context.Context, QueueUrl string) error {
	if c.opts.maxReceiveCount > 0 && len(c.opts.deadLetterQueue) > 0 {
		policy := map[string]string{
			"deadLetterTargetArn": c.opts.deadLetterQueue,
			"maxReceiveCount":     strconv.Itoa(c.opts.maxReceiveCount),
		}
		b, err := json.Marshal(policy)
		if err != nil {
			return fmt.Errorf("failed to marshal policy on err :%s", err.Error())
		}

		_, err = c.client.SetQueueAttributesWithContext(ctx, &sqs.SetQueueAttributesInput{
			Attributes: map[string]*string{
				sqs.QueueAttributeNameRedrivePolicy: aws.String(string(b)),
			},
			QueueUrl: aws.String(QueueUrl),
		})
		if err != nil {
			return fmt.Errorf("failed to SetQueueAttributesWithContext err :%s", err.Error())
		}
		return nil
	}
	return fmt.Errorf("failed to SetQueueAttributesWithContext need to verify max_receive and dead_letter exists")
}

func (c *Client) Stop() error {
	return nil
}
