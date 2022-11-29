package crate

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/logger"
	"github.com/kubemq-io/kubemq-targets/types"
	_ "github.com/lib/pq"
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
	c.db, err = sql.Open("postgres", c.opts.connection)
	if err != nil {
		return err
	}
	err = c.db.PingContext(ctx)
	if err != nil {
		_ = c.db.Close()
		return fmt.Errorf("error reaching crate at %s: %w", c.opts.connection, err)
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
	}

	return nil, nil
}

func getStatements(data []byte) []string {
	if data == nil {
		return nil
	}
	return strings.Split(string(data), ";")
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
	var results []map[string]interface{}
	for rows.Next() {
		results = append(results, parseToMap(rows, cols))
	}
	if results == nil {
		return nil
	}
	b, _ := json.Marshal(results)
	return b
}

func parseToMap(rows *sql.Rows, cols []string) map[string]interface{} {
	values := make([]interface{}, len(cols))
	pointers := make([]interface{}, len(cols))
	for i := range values {
		pointers[i] = &values[i]
	}

	if err := rows.Scan(pointers...); err != nil {
		return nil
	}

	m := make(map[string]interface{})
	for i, colName := range cols {
		if values[i] == nil {
			// m[colName] = nil
		} else {
			m[colName] = values[i]
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
