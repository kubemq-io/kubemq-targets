package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Azure/azure-storage-queue-go/azqueue"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"net/url"
)

type Client struct {
	name        string
	opts        options
	retryOption azqueue.RetryOptions
	credential  *azqueue.SharedKeyCredential
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
	// Create a default request pipeline using your storage account name and account key.
	c.credential, err = azqueue.NewSharedKeyCredential(c.opts.storageAccount, c.opts.storageAccessKey)
	if err != nil {
		return fmt.Errorf("failed to create shared key credential on error %s , please check storage access key and acccount are correct", err.Error())
	}

	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "create":
		return c.create(ctx, meta)
	case "delete":
		return c.delete(ctx, meta)
	case "get_messages_count":
		return c.getMessageCount(ctx, meta)
	case "peek":
		return c.peek(ctx, meta)
	case "push":
		return c.push(ctx, meta, req.Data)
	case "pop":
		return c.pop(ctx, meta)
	}
	return nil, errors.New("invalid method type")
}

func (c *Client) create(ctx context.Context, meta metadata) (*types.Response, error) {

	url, err := url.Parse(fmt.Sprintf("%s/%s", meta.serviceUrl, meta.queueName))
	if err != nil {
		return nil, err
	}
	queueUrl := azqueue.NewQueueURL(*url, azqueue.NewPipeline(c.credential, azqueue.PipelineOptions{
		Retry: c.retryOption,
	}))
	_, err = queueUrl.GetProperties(ctx)
	if err != nil {
		// https://godoc.org/github.com/Azure/azure-storage-queue-go/azqueue#StorageErrorCodeType
		errorType := err.(azqueue.StorageError).ServiceCode()

		if errorType == azqueue.ServiceCodeQueueNotFound {
			if len(meta.queueMetadata) > 0 {
				_, err = queueUrl.Create(ctx, meta.queueMetadata)
				if err != nil {
					return nil, err
				}
			} else {
				_, err = queueUrl.Create(ctx, azqueue.Metadata{})
				if err != nil {
					return nil, err
				}
			}
			_, err := queueUrl.GetProperties(ctx)
			if err != nil {
				return nil, err
			}

		} else {
			return nil, err
		}
	} else {
		return nil, errors.New("queue already exists")
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) delete(ctx context.Context, meta metadata) (*types.Response, error) {

	url, err := url.Parse(fmt.Sprintf("%s/%s", meta.serviceUrl, meta.queueName))
	if err != nil {
		return nil, err
	}
	queueUrl := azqueue.NewQueueURL(*url, azqueue.NewPipeline(c.credential, azqueue.PipelineOptions{
		Retry: c.retryOption,
	}))
	_, err = queueUrl.GetProperties(ctx)
	if err != nil {
		return nil, err
	}
	_, err = queueUrl.Delete(ctx)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) getMessageCount(ctx context.Context, meta metadata) (*types.Response, error) {

	url, err := url.Parse(fmt.Sprintf("%s/%s", meta.serviceUrl, meta.queueName))
	if err != nil {
		return nil, err
	}
	queueUrl := azqueue.NewQueueURL(*url, azqueue.NewPipeline(c.credential, azqueue.PipelineOptions{
		Retry: c.retryOption,
	}))
	props, err := queueUrl.GetProperties(ctx)
	if err != nil {
		return nil, err
	}
	messageCount := fmt.Sprintf("%v", props.ApproximateMessagesCount())

	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetMetadataKeyValue("count", messageCount),
		nil
}

func (c *Client) peek(ctx context.Context, meta metadata) (*types.Response, error) {

	url, err := url.Parse(fmt.Sprintf("%s/%s", meta.serviceUrl, meta.queueName))
	if err != nil {
		return nil, err
	}
	queueUrl := azqueue.NewQueueURL(*url, azqueue.NewPipeline(c.credential, azqueue.PipelineOptions{
		Retry: c.retryOption,
	}))
	messageUrl := queueUrl.NewMessagesURL()
	resp, err := messageUrl.Peek(ctx, meta.maxMessages)
	if err != nil {
		return nil, err
	}
	messages := make([]*azqueue.PeekedMessage, 0)
	for i := int32(0); i < resp.NumMessages(); i++ {
		msg := resp.Message(i)
		messages = append(messages, msg)
	}
	b, err := json.Marshal(messages)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetData(b).
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) push(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {

	url, err := url.Parse(fmt.Sprintf("%s/%s", meta.serviceUrl, meta.queueName))
	if err != nil {
		return nil, err
	}
	queueUrl := azqueue.NewQueueURL(*url, azqueue.NewPipeline(c.credential, azqueue.PipelineOptions{
		Retry: c.retryOption,
	}))
	var message string
	err = json.Unmarshal(data, &message)
	if err != nil {
		return nil, err
	}
	messageUrl := queueUrl.NewMessagesURL()
	_, err = messageUrl.Enqueue(ctx, message, meta.visibilityTimeout, meta.timeToLive)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) pop(ctx context.Context, meta metadata) (*types.Response, error) {

	url, err := url.Parse(fmt.Sprintf("%s/%s", meta.serviceUrl, meta.queueName))
	if err != nil {
		return nil, err
	}
	queueUrl := azqueue.NewQueueURL(*url, azqueue.NewPipeline(c.credential, azqueue.PipelineOptions{
		Retry: c.retryOption,
	}))
	messageUrl := queueUrl.NewMessagesURL()
	resp, err := messageUrl.Dequeue(ctx, meta.maxMessages, meta.visibilityTimeout)
	if err != nil {
		return nil, err
	}
	messages := make([]*azqueue.DequeuedMessage, 0)
	for i := int32(0); i < resp.NumMessages(); i++ {
		msg := resp.Message(i)
		messages = append(messages, msg)
		msgIdUrl := messageUrl.NewMessageIDURL(msg.ID)

		// PopReciept is required to delete the Message. If deletion fails using this popreceipt then the message has
		_, err = msgIdUrl.Delete(ctx, msg.PopReceipt)
		if err != nil {
			return nil, err
		}
	}
	b, err := json.Marshal(messages)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetData(b).
			SetMetadataKeyValue("result", "ok"),
		nil
}
