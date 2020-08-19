package athena

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type Client struct {
	name   string
	opts   options
	client *athena.Athena
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

	svc := athena.New(sess)
	c.client = svc

	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "list_databases":
		return c.listDatabases(ctx, meta)
	case "list_data_catalogs":
		return c.listDataCatalogs(ctx)
	case "query":
		return c.query(ctx, meta)
	case "get_query_result":
		return c.getQueryResult(ctx, meta)
	default:
		return nil, fmt.Errorf("invalid method type")
	}
}

func (c *Client) listDatabases(ctx context.Context, meta metadata) (*types.Response, error) {
	m, err := c.client.ListDatabasesWithContext(ctx, &athena.ListDatabasesInput{
		CatalogName: aws.String(meta.catalog),
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

func (c *Client) listDataCatalogs(ctx context.Context) (*types.Response, error) {
	m, err := c.client.ListDataCatalogsWithContext(ctx, nil)
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

func (c *Client) query(ctx context.Context, meta metadata) (*types.Response, error) {
	m, err := c.client.StartQueryExecutionWithContext(ctx, &athena.StartQueryExecutionInput{
		QueryString: aws.String(meta.query),
		QueryExecutionContext: &athena.QueryExecutionContext{
			Catalog:  aws.String(meta.catalog),
			Database: aws.String(meta.DB),
		},
		ResultConfiguration: &athena.ResultConfiguration{
			OutputLocation: aws.String(meta.outputLocation),
		},
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
			SetMetadataKeyValue("execution_id", *m.QueryExecutionId).
			SetData(b),
		nil
}

func (c *Client) getQueryResult(ctx context.Context, meta metadata) (*types.Response, error) {
	m, err := c.client.GetQueryResultsWithContext(ctx, &athena.GetQueryResultsInput{
		QueryExecutionId: aws.String(meta.executionID),
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
