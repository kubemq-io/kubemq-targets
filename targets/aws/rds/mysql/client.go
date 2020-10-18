package mysql

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/go-sql-driver/mysql"
	"github.com/kubemq-hub/builder/connector/common"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	_ "github.com/go-sql-driver/mysql"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Client is a Client state store
type Client struct {
	name string
	db   *sql.DB
	opts options
}

func New() *Client {
	return &Client{}
}
func (c *Client) Connector() *common.Connector {
	return Connector()
}
func (c *Client) Init(ctx context.Context, cfg config.Spec) error {

	c.name = cfg.Name
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}

	host := fmt.Sprintf("%s:%d", c.opts.endPoint, c.opts.dbPort)
	mysqlCfp := &mysql.Config{
		User: c.opts.dbUser,
		Addr: host,
		Net:  "tcp",
		Params: map[string]string{
			"tls": "rds",
		},
		DBName: c.opts.dbName,
	}
	mysqlCfp.AllowNativePasswords = true
	mysqlCfp.AllowCleartextPasswords = true

	mysqlCfp.Passwd, err = rdsutils.BuildAuthToken(host, c.opts.region, mysqlCfp.User, credentials.NewStaticCredentials(c.opts.awsKey, c.opts.awsSecretKey, c.opts.token))
	if err != nil {
		return err
	}

	err = registerRDSMysqlCerts(http.DefaultClient)
	if err != nil {
		return err
	}
	a := mysqlCfp.FormatDSN()
	c.db, err = sql.Open("mysql", a)
	if err != nil {
		return err
	}
	err = c.db.PingContext(ctx)
	if err != nil {
		return err
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

	return nil, errors.New("invalid method type")
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
func (c *Client) CloseClient() error {
	return c.db.Close()
}

//https://github.com/aws/aws-sdk-go/issues/1248
func registerRDSMysqlCerts(c *http.Client) error {
	resp, err := c.Get("https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem")
	if err != nil {
		return err
	}
	pem, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return err
	}

	err = mysql.RegisterTLSConfig("rds", &tls.Config{RootCAs: rootCertPool})
	if err != nil {
		return err
	}
	return nil
}
