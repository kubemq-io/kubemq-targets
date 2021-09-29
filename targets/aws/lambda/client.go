package lambda

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
)

type Client struct {
	log    *logger.Logger
	opts   options
	client *lambda.Lambda
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

	svc := lambda.New(sess)
	c.client = svc

	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "list":
		return c.list(ctx)
	case "create":
		return c.create(ctx, meta, req.Data)
	case "run":
		return c.run(ctx, meta, req.Data)
	case "delete":
		return c.delete(ctx, meta)
	default:
		return nil, errors.New("invalid method type")
	}
}

func (c *Client) list(ctx context.Context) (*types.Response, error) {
	m, err := c.client.ListFunctionsWithContext(ctx, nil)
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

func (c *Client) create(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	if data == nil {
		return nil, fmt.Errorf("data is empty , please add lambda body as []byte")
	}

	input := &lambda.CreateFunctionInput{
		Code: &lambda.FunctionCode{
			ZipFile: data,
		},
		Description:  aws.String(meta.description),
		FunctionName: aws.String(meta.functionName),
		Handler:      aws.String(meta.handlerName),
		MemorySize:   aws.Int64(meta.memorySize),
		Role:         aws.String(meta.role),
		Runtime:      aws.String(meta.runtime),
		Timeout:      aws.Int64(meta.timeout),
	}

	result, err := c.client.CreateFunctionWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) run(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	if isValid := isJson(data); !isValid {
		return nil, fmt.Errorf("function payload must be a valid json")
	}
	invoke := &lambda.InvokeInput{}
	invoke.SetFunctionName(meta.functionName).SetPayload(data)
	if meta.isDryRun {
		invoke.SetInvocationType("DryRun")
	}
	result, err := c.client.InvokeWithContext(ctx, invoke)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) delete(ctx context.Context, meta metadata) (*types.Response, error) {

	_, err := c.client.DeleteFunctionWithContext(ctx, &lambda.DeleteFunctionInput{FunctionName: aws.String(meta.functionName)})
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}
func isJson(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	obj := map[string]interface{}{}
	err := json.Unmarshal(data, &obj)
	return err == nil
}

func (c *Client) Stop() error {
	return nil
}
