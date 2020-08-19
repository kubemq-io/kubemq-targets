package dynamodb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type Client struct {
	name   string
	opts   options
	client *dynamodb.DynamoDB
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

	svc := dynamodb.New(sess)
	c.client = svc

	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "list_tables":
		return c.listTables(ctx)
	case "create_table":
		return c.createTable(ctx, req.Data)
	case "delete_table":
		return c.deleteTable(ctx, meta)
	case "insert_item":
		return c.insertItem(ctx, meta, req.Data)
	case "get_item":
		return c.getItem(ctx, req.Data)
	case "update_item":
		return c.updateItem(ctx, req.Data)
	case "delete_item":
		return c.deleteItem(ctx, req.Data)
	default:
		return nil, fmt.Errorf(getValidMethodTypes())
	}
}

func (c *Client) listTables(ctx context.Context) (*types.Response, error) {
	input := &dynamodb.ListTablesInput{}
	m, err := c.client.ListTablesWithContext(ctx, input)
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

func (c *Client) createTable(ctx context.Context, data []byte) (*types.Response, error) {
	i := &dynamodb.CreateTableInput{}
	err := json.Unmarshal(data, &i)
	if err != nil {
		return nil, err
	}

	result, err := c.client.CreateTableWithContext(ctx, i)
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

func (c *Client) deleteTable(ctx context.Context, meta metadata) (*types.Response, error) {
	result, err := c.client.DeleteTableWithContext(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(meta.tableName),
	})
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
func (c *Client) insertItem(ctx context.Context, meta metadata, data []byte) (*types.Response, error) {
	i := map[string]*dynamodb.AttributeValue{}
	err := json.Unmarshal(data, &i)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      i,
		TableName: aws.String(meta.tableName),
	}
	result, err := c.client.PutItemWithContext(ctx, input)
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

func (c *Client) getItem(ctx context.Context, data []byte) (*types.Response, error) {
	g := &dynamodb.GetItemInput{}
	err := json.Unmarshal(data, &g)
	if err != nil {
		return nil, err
	}
	result, err := c.client.GetItemWithContext(ctx, g)
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


func (c *Client) updateItem(ctx context.Context, data []byte) (*types.Response, error) {
	u := &dynamodb.UpdateItemInput{}
	err := json.Unmarshal(data, &u)
	if err != nil {
		return nil, err
	}
	result, err := c.client.UpdateItemWithContext(ctx, u)
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


func (c *Client) deleteItem(ctx context.Context, data []byte) (*types.Response, error) {
	d := &dynamodb.DeleteItemInput{}
	err := json.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}
	result, err := c.client.DeleteItemWithContext(ctx, d)
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