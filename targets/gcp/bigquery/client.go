package bigquery

import (
	"cloud.google.com/go/bigquery"
	"context"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
	"github.com/kubemq-hub/kubemq-targets/types"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Client struct {
	opts   options
	client *bigquery.Client
	log    *logger.Logger
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

	c.log = logger.NewLogger(cfg.Name)
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	b := []byte(c.opts.credentials)
	Client, err := bigquery.NewClient(ctx, c.opts.projectID, option.WithCredentialsJSON(b))
	if err != nil {
		return err
	}
	c.client = Client
	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "query":
		return c.query(ctx, meta)
	case "create_table":
		return c.createTable(ctx, meta, req.Data)
	case "delete_table":
		return c.deleteTable(ctx, meta)
	case "create_data_set":
		return c.createDataSet(ctx, meta)
	case "delete_data_set":
		return c.deleteDataSet(ctx, meta)
	case "get_table_info":
		return c.getTableInfo(ctx, meta)
	case "get_data_sets":
		return c.getDataSets(ctx)
	case "insert":
		return c.insert(ctx, meta, req.Data)
	default:
		return nil, errors.New("invalid method type")
	}
}

func (c *Client) getTableInfo(ctx context.Context, meta metadata) (*types.Response, error) {
	m, err := c.client.Dataset(meta.datasetID).Table(meta.tableName).Metadata(ctx)
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

func (c *Client) getDataSets(ctx context.Context) (*types.Response, error) {
	i := c.client.Datasets(ctx)
	s, err := c.getDataSetsFromIterator(i)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	if len(s) == 0 {
		return nil, fmt.Errorf("no data sets found")
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) query(ctx context.Context, meta metadata) (*types.Response, error) {
	query := c.client.Query(meta.query)
	i, err := query.Read(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := c.getRowsFromIterator(i)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows found")
	}
	b, err := json.Marshal(rows)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) createDataSet(ctx context.Context, meta metadata) (*types.Response, error) {
	met := &bigquery.DatasetMetadata{
		Location: meta.location, // See https://cloud.google.com/bigquery/docs/locations
	}
	err := c.client.Dataset(meta.datasetID).Create(ctx, met)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) deleteDataSet(ctx context.Context, meta metadata) (*types.Response, error) {
	err := c.client.Dataset(meta.datasetID).Delete(ctx)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) insert(ctx context.Context, meta metadata, body []byte) (*types.Response, error) {
	ir, err := newInsertRecord(body)
	if err != nil {
		return nil, err
	}
	ins := c.client.Dataset(meta.datasetID).Table(meta.tableName).Inserter()
	err = ins.Put(ctx, ir.records)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetMetadataKeyValue("insert_rows", fmt.Sprintf("%d", len(ir.records))),
		nil
}

func (c *Client) createTable(ctx context.Context, meta metadata, body []byte) (*types.Response, error) {
	metaData := &bigquery.TableMetadata{}

	err := json.Unmarshal(body, &metaData)
	if err != nil {
		return nil, err
	}
	tableRef := c.client.Dataset(meta.datasetID).Table(meta.tableName)
	err = tableRef.Create(ctx, metaData)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) deleteTable(ctx context.Context, meta metadata) (*types.Response, error) {
	tableRef := c.client.Dataset(meta.datasetID).Table(meta.tableName)
	err := tableRef.Delete(ctx)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) getRowsFromIterator(iter *bigquery.RowIterator) ([]map[string]bigquery.Value, error) {
	var rows []map[string]bigquery.Value
	for {
		row := make(map[string]bigquery.Value)
		err := iter.Next(&row)
		if err == iterator.Done {
			return rows, nil
		}
		if err != nil {
			return rows, fmt.Errorf("error iterating through results: %v", err)
		}
		rows = append(rows, row)
	}
}

func (c *Client) getDataSetsFromIterator(iter *bigquery.DatasetIterator) ([]*bigquery.Dataset, error) {
	var datasets []*bigquery.Dataset
	for {
		dataset, err := iter.Next()
		if err == iterator.Done {
			return datasets, nil
		}
		if err != nil {
			return datasets, fmt.Errorf("error iterating through results: %v", err)
		}
		datasets = append(datasets, dataset)
	}
}

func (c *Client) Stop() error {
	return c.client.Close()
}
