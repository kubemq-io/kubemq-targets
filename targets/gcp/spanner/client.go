package spanner

import (
	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/pkg/logger"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

type Client struct {
	name        string
	opts        options
	client      *spanner.Client
	adminClient *database.DatabaseAdminClient
	log         *logger.Logger
}

func New() *Client {
	return &Client{}

}
func (c *Client) Name() string {
	return c.name
}

func (c *Client) Init(ctx context.Context, cfg config.Metadata) error {
	c.name = cfg.Name
	c.log = logger.NewLogger(cfg.Name)
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	b := []byte(c.opts.credentials)
	adminClient, err := database.NewDatabaseAdminClient(ctx,option.WithCredentialsJSON(b))
	if err != nil {
		return err
	}

	c.adminClient = adminClient
	Client, err := spanner.NewClient(ctx, c.opts.db,option.WithCredentialsJSON(b))
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
	case "read":
		return c.read(ctx, meta, req.Data)
	case "update_database_ddl":
		return c.updateDatabaseDdl(ctx, req.Data)
	case "insert":
		return c.insert(ctx, req.Data)
	case "update":
		return c.update(ctx, req.Data)
	case "insert_or_update":
		return c.insertOrUpdate(ctx, req.Data)
	default:
		return nil, fmt.Errorf(getValidMethodTypes())
	}
}

func (c *Client) query(ctx context.Context, meta metadata) (*types.Response, error) {
	stmt := spanner.Statement{SQL: meta.query}
	i := c.client.Single().Query(ctx, stmt)
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
	var q []spanner.Row
	err = json.Unmarshal(b, &q)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) insert(ctx context.Context, body []byte) (*types.Response, error) {
	var inserts []InsertOrUpdate
	err := json.Unmarshal(body, &inserts)
	if err != nil {
		return nil, err
	}
	if len(inserts) == 0 {
		return nil, fmt.Errorf("failed to get valid InsertOrUpdate struct")
	}
	var m []*spanner.Mutation
	for _, i := range inserts {
		err = i.validate()
		if err != nil {
			return nil, err
		}
		m = append(m, spanner.Insert(i.TableName, i.ColumnName, i.ColumnValue))
	}
	_, err = c.client.Apply(ctx, m)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) update(ctx context.Context, body []byte) (*types.Response, error) {
	var updates []InsertOrUpdate
	err := json.Unmarshal(body, &updates)
	if err != nil {
		return nil, err
	}
	if len(updates) == 0 {
		return nil, fmt.Errorf("failed to get valid InsertOrUpdate struct")
	}
	var m []*spanner.Mutation
	for _, i := range updates {
		err = i.validate()
		if err != nil {
			return nil, err
		}
		m = append(m, spanner.Update(i.TableName, i.ColumnName, i.ColumnValue))
	}
	_, err = c.client.Apply(ctx, m)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) insertOrUpdate(ctx context.Context, body []byte) (*types.Response, error) {
	var insertsOrUpdates []InsertOrUpdate
	err := json.Unmarshal(body, &insertsOrUpdates)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if len(insertsOrUpdates) == 0 {
		return nil, fmt.Errorf("failed to get valid InsertOrUpdate struct")
	}
	var m []*spanner.Mutation
	for _, i := range insertsOrUpdates {
		err = i.validate()
		if err != nil {
			return nil, err
		}
		m = append(m, spanner.InsertOrUpdate(i.TableName, i.ColumnName, i.ColumnValue))
	}
	_, err = c.client.Apply(ctx, m)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) read(ctx context.Context, meta metadata, body []byte) (*types.Response, error) {
	var columns []string
	err := json.Unmarshal(body, &columns)
	if err != nil {
		return nil, fmt.Errorf("failed to parse body as []strings for columns on error %s", err.Error())
	}
	iter := c.client.Single().Read(ctx, meta.tableName, spanner.AllKeys(), columns)
	rows, err := c.getRowsFromIterator(iter)
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

func (c *Client) getRowsFromIterator(iter *spanner.RowIterator) ([]*Row, error) {
	var rows []*Row
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return rows, nil
		}
		if err != nil {
			return rows, fmt.Errorf("error iterating through results: %v", err)
		}
		r, err := extractDataByType(row)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r)
	}
}

func (c *Client) updateDatabaseDdl(ctx context.Context, body []byte) (*types.Response, error) {
	var Statements []string
	err := json.Unmarshal(body, &Statements)
	if err != nil {
		return nil, fmt.Errorf("failed to parse body as []strings for args on error %s", err.Error())
	}
	if len(Statements) == 0 {
		return nil, fmt.Errorf("failed to get valid Statements struct")
	}
	op, err := c.adminClient.UpdateDatabaseDdl(ctx, &adminpb.UpdateDatabaseDdlRequest{
		Database:   c.opts.db,
		Statements: Statements,
	})
	if err != nil {
		return nil, err
	}
	// nolint:staticcheck
	b, err := json.Marshal(op)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok").
			SetData(b),
		nil
}

func (c *Client) CloseClient() error {
	c.client.Close()
	return c.adminClient.Close()
}
