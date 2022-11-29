package metrics

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
)

type Client struct {
	log    *logger.Logger
	opts   options
	client *cloudwatch.CloudWatch
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

	svc := cloudwatch.New(sess)
	c.client = svc

	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "put_metrics":
		return c.putMetrics(ctx, meta, req.Data)
	case "list_metrics":
		return c.listMetrics(ctx, meta)
	default:
		return nil, errors.New("invalid method type")
	}
}

func (c *Client) putMetrics(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	var metrics []*cloudwatch.MetricDatum
	err := json.Unmarshal(data, &metrics)
	if err != nil {
		return nil, err
	}
	_, err = c.client.PutMetricDataWithContext(ctx, &cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(meta.namespace),
		MetricData: metrics,
	})
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) listMetrics(ctx context.Context, meta metadata) (*types.Response, error) {
	var resp *cloudwatch.ListMetricsOutput
	var err error
	if meta.namespace != "" {
		resp, err = c.client.ListMetricsWithContext(ctx, &cloudwatch.ListMetricsInput{
			Namespace: aws.String(meta.namespace),
		})
	} else {
		resp, err = c.client.ListMetricsWithContext(ctx, &cloudwatch.ListMetricsInput{})
	}
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

func (c *Client) Stop() error {
	return nil
}
