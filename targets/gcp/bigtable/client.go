package bigtable

import (
	"bytes"
	"cloud.google.com/go/bigtable"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
)

type Client struct {
	name        string
	opts        options
	adminClient *bigtable.AdminClient
	client      *bigtable.Client
}

func New() *Client {
	return &Client{}

}
func (c *Client) Name() string {
	return c.name
}

func (c *Client) Init(ctx context.Context, cfg config.Metadata) error {
	c.name = cfg.Name
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}

	adminClient, err := bigtable.NewAdminClient(ctx, c.opts.projectID, c.opts.instance)
	if err != nil {
		return err
	}
	c.adminClient = adminClient

	Client, err := bigtable.NewClient(ctx, c.opts.projectID, c.opts.instance)
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
	case "write":
		return c.writeRow(ctx, meta, req.Data)
	case "write_batch":
		return c.writeBatch(ctx, meta, req.Data)
	case "get_row":
		return c.readRow(ctx, meta)
	case "get_all_rows":
		return c.readAllRows(ctx, meta, req.Data)
	case "get_all_rows_by_column":
		return c.readAllRowsByColumnFilter(ctx, meta, req.Data)
	case "create_column_family":
		return c.createColumnFamily(ctx, meta)
	case "delete_row":
		return c.deleteRowRange(ctx, meta)
	case "get_tables":
		return c.getTables(ctx)
	case "create_table":
		return c.createTable(ctx, meta)
	case "delete_table":
		return c.deleteTable(ctx, meta)
	default:
		return nil, fmt.Errorf(getValidMethodTypes())
	}
}

func (c *Client) getTables(ctx context.Context) (*types.Response, error) {
	tables, err := c.adminClient.Tables(ctx)
	if err!= nil {
		return nil, err
	}
	if len(tables) <= 0 {
		return nil, fmt.Errorf("no tables found for this instance")

	}
	b, err := json.Marshal(tables)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetData(b).
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) createTable(ctx context.Context, meta metadata) (*types.Response, error) {
	err := c.adminClient.CreateTable(ctx, meta.tableName)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) createColumnFamily(ctx context.Context, meta metadata) (*types.Response, error) {
	err := c.adminClient.CreateColumnFamily(ctx, meta.tableName, meta.columnFamily)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) deleteTable(ctx context.Context, meta metadata) (*types.Response, error) {
	err := c.adminClient.DeleteTable(ctx, meta.tableName)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) deleteRowRange(ctx context.Context, meta metadata) (*types.Response, error) {
	err := c.adminClient.DropRowRange(ctx, meta.tableName, meta.rowKeyPrefix)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) readRow(ctx context.Context, meta metadata) (*types.Response, error) {
	tbl := c.client.Open(meta.tableName)
	defer func() {
		_=c.client.Close()
	}()
	row, err := tbl.ReadRow(ctx, meta.rowKeyPrefix)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(row)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetData(b),
		nil
}

func (c *Client) readAllRows(ctx context.Context, meta metadata, body []byte) (*types.Response, error) {
	prefixes, err := c.getRowKeys(body)
	if err != nil {
		return nil, err
	}
	var rs bigtable.RowSet
	if len(prefixes) > 0 {
		rs = bigtable.RowList(prefixes)
	} else {
		rs = bigtable.PrefixRange("")
	}
	tbl := c.client.Open(meta.tableName)
	defer c.client.Close()
	var rows []bigtable.Row
	err = tbl.ReadRows(ctx, rs, func(row bigtable.Row) bool {
		rows = append(rows, row)
		return true
	})
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows found for this table")
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

func (c *Client) readAllRowsByColumnFilter(ctx context.Context, meta metadata, body []byte) (*types.Response, error) {
	prefixes, err := c.getRowKeys(body)
	if err != nil {
		return nil, err
	}
	var rs bigtable.RowSet
	if len(prefixes) > 0 {
		rs = bigtable.RowList(prefixes)
	} else {
		rs = bigtable.PrefixRange("")
	}
	tbl := c.client.Open(meta.tableName)
	defer c.client.Close()
	var rows []bigtable.Row
	err = tbl.ReadRows(ctx, rs, func(row bigtable.Row) bool {
		rows = append(rows, row)
		return true
	}, bigtable.RowFilter(bigtable.ColumnFilter(meta.readColumnName)))
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows found for this table")
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

func (c *Client) writeRow(ctx context.Context, meta metadata, body []byte) (*types.Response, error) {
	tbl := c.client.Open(meta.tableName)
	defer func() {
		_=c.client.Close()
	}()
	timestamp := bigtable.Now()
	mut := bigtable.NewMutation()
	m, err := c.getSingleColumnFromBody(body)
	if err!=nil {
		return nil, err
	}
	rowKey := ""
	for k, v := range m {
		if k == "set_row_key" {
			rowKey = fmt.Sprintf("%s", v)
		} else {
			buf := new(bytes.Buffer)
			b, err := json.Marshal(v)
			if err != nil {
				return types.NewResponse().
					SetMetadataKeyValue("error", "true").
					SetMetadataKeyValue("message", err.Error()), nil
			}
			err = binary.Write(buf, binary.BigEndian, b)
			if err != nil {
				return types.NewResponse().
					SetMetadataKeyValue("error", "true").
					SetMetadataKeyValue("message", err.Error()), nil
			}
			mut.Set(meta.columnFamily, k, timestamp, buf.Bytes())
		}
	}
	if len(rowKey) == 0 {
		return nil, fmt.Errorf( "missing set_row_key value")
	}
	err = tbl.Apply(ctx, rowKey, mut)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) writeBatch(ctx context.Context, meta metadata, body []byte) (*types.Response, error) {
	tbl := c.client.Open(meta.tableName)
	defer func() {
		_=c.client.Close()
	}()
	timestamp := bigtable.Now()
	var muts []*bigtable.Mutation
	var rowKeys []string
	s, err := c.getMultipleColumnsFromBody(body)
	if err!=nil {
		return nil, err
	}
	if len(s) == 0 {
		return nil,fmt.Errorf("column requested must be at least 1")
	}
	for _, m := range s {
		mut := bigtable.NewMutation()
		for k, v := range m {
			if k == "set_row_key" {
				rowKeys = append(rowKeys, fmt.Sprintf("%s", v))
			} else {
				buf := new(bytes.Buffer)
				b, err := json.Marshal(v)
				if err != nil {
					return types.NewResponse().
						SetMetadataKeyValue("error", "true").
						SetMetadataKeyValue("message", err.Error()), nil
				}
				_=binary.Write(buf, binary.BigEndian, b)
				mut.Set(meta.columnFamily, k, timestamp, buf.Bytes())
			}
		}
		muts = append(muts, mut)
	}
	if len(s) != len(rowKeys) {
		return nil,fmt.Errorf("set_row_key count does not match column requested")
	}
	_, err = tbl.ApplyBulk(ctx, rowKeys, muts)
	if err != nil {
		return nil,err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) getSingleColumnFromBody(body []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(body, &m)
	if err != nil {
		return nil, fmt.Errorf("failed to parse body as map[string]interface{} on error %s", err.Error())
	}
	return m, nil
}

func (c *Client) getRowKeys(body []byte) ([]string, error) {
	var m []string
	if body != nil {
		err := json.Unmarshal(body, &m)
		if err != nil {
			return nil, fmt.Errorf("failed to parse body as []string on error %s", err.Error())
		}
	}
	return m, nil
}

func (c *Client) getMultipleColumnsFromBody(body []byte) ([]map[string]interface{}, error) {
	var s []map[string]interface{}
	err := json.Unmarshal(body, &s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse body as slice of map[string]interface{} on error %s", err.Error())
	}
	return s, nil
}

func (c *Client) CloseAdminClient() error {
	return c.adminClient.Close()
}
