package percona

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
	"github.com/kubemq-hub/kubemq-targets/types"
	"strconv"
	"strings"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Client is a Client state store
type Client struct {
	log  *logger.Logger
	db   *sql.DB
	opts options
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
	c.db, err = sql.Open("mysql", c.opts.connection)
	if err != nil {
		return err
	}
	err = c.db.PingContext(ctx)
	if err != nil {
		_ = c.db.Close()
		return fmt.Errorf("error reaching mysql at %s: %w", c.opts.connection, err)
	}
	c.db.SetMaxOpenConns(c.opts.maxOpenConnections)
	c.db.SetMaxIdleConns(c.opts.maxIdleConnections)
	c.db.SetConnMaxLifetime(time.Duration(c.opts.connectionMaxLifetimeSeconds) * time.Second)
	return nil
}

func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "query":
		return c.Query(ctx, meta, req.Data)
	case "exec":
		return c.Exec(ctx, meta, req.Data)
	case "transaction":
		return c.Transaction(ctx, meta, req.Data)
	}

	return nil, nil
}
func (c *Client) Exec(ctx context.Context, meta metadata, value []byte) (*types.Response, error) {
	stmts := getStatements(value)
	if stmts == nil {
		return nil, fmt.Errorf("no exec statement found")
	}
	for i, stmt := range stmts {
		if stmt != "" {
			_, err := c.db.ExecContext(ctx, stmt)
			if err != nil {
				return nil, fmt.Errorf("error on statement %d, %w", i, err)
			}
		}
	}

	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}
func getStatements(data []byte) []string {
	if data == nil {
		return nil
	}
	return strings.Split(string(data), ";")
}
func (c *Client) Transaction(ctx context.Context, meta metadata, value []byte) (*types.Response, error) {
	stmts := getStatements(value)
	if stmts == nil {
		return nil, fmt.Errorf("no transaction statements found")
	}

	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: meta.isolationLevel,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()
	for i, stmt := range stmts {
		if stmt != "" {
			_, err := tx.ExecContext(ctx, stmt)
			if err != nil {
				rollBackErr := tx.Rollback()
				if rollBackErr != nil {
					return nil, rollBackErr
				}
				return nil, fmt.Errorf("error on statement %d, %w", i, err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) Query(ctx context.Context, meta metadata, value []byte) (*types.Response, error) {
	stmt := string(value)
	if stmt == "" {
		return nil, fmt.Errorf("no query statement found")
	}
	rows, err := c.db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return types.NewResponse().
		SetData(c.rowsToMap(rows)).
		SetMetadataKeyValue("result", "ok"), nil

}

func (c *Client) rowsToMap(rows *sql.Rows) []byte {

	cols, _ := rows.Columns()
	colsTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil
	}
	var results []map[string]interface{}
	for rows.Next() {
		results = append(results, parseWithRawBytes(rows, cols, colsTypes))
	}
	if results == nil {
		return nil
	}
	b, _ := json.Marshal(results)
	return b
}

func parseWithRawBytes(rows *sql.Rows, cols []string, colsTypes []*sql.ColumnType) map[string]interface{} {
	vals := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(vals))
	for i := range vals {
		scanArgs[i] = &vals[i]
	}
	if err := rows.Scan(scanArgs...); err != nil {
		panic(err)
	}
	m := make(map[string]interface{})
	for i, col := range vals {
		if col == nil {
			continue
		}
		switch colsTypes[i].DatabaseTypeName() {
		case "TINYINT", "BOOLEAN":
			m[cols[i]], _ = strconv.ParseBool(string(col))
		case "SMALLINT", "MEDIUMINT":
			m[cols[i]], _ = strconv.Atoi(string(col))
		case "BIGINT":
			m[cols[i]], _ = strconv.ParseInt(string(col), 10, 64)
		case "FLOAT":
			val, _ := strconv.ParseFloat(string(col), 32)
			m[cols[i]] = float32(val)
		case "DOUBLE", "DECIMAL":
			m[cols[i]], _ = strconv.ParseFloat(string(col), 64)
		default:
			m[cols[i]] = string(col)
		}
	}
	return m
}

func (c *Client) Stop() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
