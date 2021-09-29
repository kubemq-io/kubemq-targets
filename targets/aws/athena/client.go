package athena

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
)

type Client struct {
	log    *logger.Logger
	opts   options
	client *athena.Athena
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
		return nil, errors.New("invalid method type")
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

func (c *Client) Stop() error {
	return nil
}
