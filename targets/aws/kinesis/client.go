package kinesis

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
)

type Client struct {
	log    *logger.Logger
	opts   options
	client *kinesis.Kinesis
}

func New() *Client {
	return &Client{}
}

func (c *Client) Connector() *common.Connector {
	return Connector()
}

func (c *Client) Init(ctx context.Context, cfg config.Spec, log *logger.Logger) error {
	c.log = log
	if c.log == nil {
		c.log = logger.NewLogger(cfg.Kind)
	}

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

	svc := kinesis.New(sess)
	c.client = svc

	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "list_streams":
		return c.listStreams(ctx)
	case "list_stream_consumers":
		return c.listStreamConsumers(ctx, meta)
	case "create_stream":
		return c.createStream(ctx, meta)
	case "create_stream_consumer":
		return c.createStreamConsumer(ctx, meta)
	case "delete_stream":
		return c.deleteStream(ctx, meta)
	case "put_record":
		return c.putRecord(ctx, meta, req.Data)
	case "put_records":
		return c.putRecords(ctx, meta, req.Data)
	case "get_records":
		return c.getRecord(ctx, meta)
	case "get_shard_iterator":
		return c.getShardIterator(ctx, meta)
	case "list_shards":
		return c.listShards(ctx, meta)
	default:
		return nil, errors.New("invalid method type")
	}
}

func (c *Client) listStreams(ctx context.Context) (*types.Response, error) {
	m, err := c.client.ListStreamsWithContext(ctx, nil)
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

func (c *Client) listStreamConsumers(ctx context.Context, meta metadata) (*types.Response, error) {
	m, err := c.client.ListStreamConsumersWithContext(ctx, &kinesis.ListStreamConsumersInput{
		StreamARN: aws.String(meta.streamARN),
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

func (c *Client) createStream(ctx context.Context, meta metadata) (*types.Response, error) {
	_, err := c.client.CreateStreamWithContext(ctx, &kinesis.CreateStreamInput{
		ShardCount: aws.Int64(meta.shardCount),
		StreamName: aws.String(meta.streamName),
	})
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) deleteStream(ctx context.Context, meta metadata) (*types.Response, error) {
	_, err := c.client.DeleteStreamWithContext(ctx, &kinesis.DeleteStreamInput{
		StreamName: aws.String(meta.streamName),
	})
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) listShards(ctx context.Context, meta metadata) (*types.Response, error) {
	m, err := c.client.ListShardsWithContext(ctx, &kinesis.ListShardsInput{
		StreamName: aws.String(meta.streamName),
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

func (c *Client) putRecord(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	m, err := c.client.PutRecordWithContext(ctx, &kinesis.PutRecordInput{
		Data:         data,
		StreamName:   aws.String(meta.streamName),
		PartitionKey: aws.String(meta.partitionKey),
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

func (c *Client) getShardIterator(ctx context.Context, meta metadata) (*types.Response, error) {
	m, err := c.client.GetShardIteratorWithContext(ctx, &kinesis.GetShardIteratorInput{
		ShardId:           aws.String(meta.shardID),
		StreamName:        aws.String(meta.streamName),
		ShardIteratorType: aws.String(meta.shardIteratorType),
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
			SetMetadataKeyValue("shard_iterator", *m.ShardIterator).
			SetData(b),
		nil
}

func (c *Client) putRecords(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	rm := make(map[string][]byte)
	err := json.Unmarshal(data, &rm)
	if err != nil {
		return nil, err
	}
	var r []*kinesis.PutRecordsRequestEntry
	for k, v := range rm {
		r = append(r, &kinesis.PutRecordsRequestEntry{Data: v, PartitionKey: aws.String(k)})
	}
	m, err := c.client.PutRecordsWithContext(ctx, &kinesis.PutRecordsInput{
		Records:    r,
		StreamName: aws.String(meta.streamName),
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

func (c *Client) getRecord(ctx context.Context, meta metadata) (*types.Response, error) {
	m, err := c.client.GetRecordsWithContext(ctx, &kinesis.GetRecordsInput{
		ShardIterator: aws.String(meta.shardIteratorID),
		Limit:         aws.Int64(meta.limit),
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

func (c *Client) createStreamConsumer(ctx context.Context, meta metadata) (*types.Response, error) {
	m, err := c.client.RegisterStreamConsumerWithContext(ctx, &kinesis.RegisterStreamConsumerInput{
		ConsumerName: aws.String(meta.consumerName),
		StreamARN:    aws.String(meta.streamARN),
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
