package keyspaces

import (
	"context"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
	"github.com/kubemq-hub/kubemq-targets/types"
	"io"
	"net/http"
	"os"
)

// Client is a Client state store
type Client struct {
	log     *logger.Logger
	session *gocql.Session
	cluster *gocql.ClusterConfig
	table   string
	opts    options
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
	c.cluster = gocql.NewCluster(c.opts.hosts...)
	err = c.downloadFile()
	if err != nil {
		return err
	}
	if c.opts.username != "" && c.opts.password != "" {
		c.cluster.Authenticator = gocql.PasswordAuthenticator{Username: c.opts.username, Password: c.opts.password}
	}
	c.cluster.Port = c.opts.port
	c.cluster.ProtoVersion = c.opts.protoVersion
	c.cluster.Consistency = c.opts.consistency
	c.cluster.Timeout = c.opts.timeoutSeconds
	c.cluster.ConnectTimeout = c.opts.connectTimeoutSeconds
	c.cluster.SslOpts = &gocql.SslOptions{
		CaPath: "tls.pem",
	}
	c.session, err = c.cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("error creating session to keyspace at %s: %w", c.opts.hosts, err)
	}
	if c.opts.defaultKeyspace != "" && c.opts.defaultTable != "" {
		err = c.tryCreateKeyspace(c.opts.defaultKeyspace, c.opts.replicationFactor)
		if err != nil {
			c.session.Close()
			return fmt.Errorf("error creating defaultKeyspace %s: %s", c.opts.defaultTable, err)
		}
		err = c.tryCreateTable(c.opts.defaultTable, c.opts.defaultKeyspace)
		if err != nil {
			c.session.Close()
			return fmt.Errorf("error creating defaultKeyspace %s: %s", c.opts.defaultTable, err)
		}
		c.table = fmt.Sprintf("%s.%s", c.opts.defaultKeyspace, c.opts.defaultTable)
	}

	return nil
}

func (c *Client) tryCreateKeyspace(keyspace string, replicationFactor int) error {
	return c.session.Query(fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : %s};", keyspace, fmt.Sprintf("%v", replicationFactor))).Exec()
}

func (c *Client) tryCreateTable(table, keyspace string) error {
	return c.session.Query(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (key text, value blob, PRIMARY KEY (key));", keyspace, table)).Exec()
}
func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "get":
		return c.Get(ctx, meta)
	case "set":
		return c.Set(ctx, meta, req.Data)
	case "delete":
		return c.Delete(ctx, meta)
	case "query":
		return c.Query(ctx, meta, req.Data)
	case "exec":
		return c.Exec(ctx, meta, req.Data)
	}
	return nil, nil
}
func (c *Client) createSession(consistency gocql.Consistency) (*gocql.Session, error) {
	session, err := c.cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("error creating session: %s", err)
	}

	session.SetConsistency(consistency)
	return session, nil
}
func (c *Client) Get(ctx context.Context, meta metadata) (*types.Response, error) {
	session := c.session

	switch meta.consistency {
	case "One":
		sess, err := c.createSession(gocql.One)
		if err != nil {
			return nil, err
		}
		defer sess.Close()
		session = sess
	case "LocalOne":
		sess, err := c.createSession(gocql.LocalOne)
		if err != nil {
			return nil, err
		}
		defer sess.Close()
		session = sess
	case "LocalQuorum":
		sess, err := c.createSession(gocql.LocalQuorum)
		if err != nil {
			return nil, err
		}
		defer sess.Close()
		session = sess
	}
	table := meta.keyspaceTable()
	if table == "" {
		table = c.table
	}
	/* #nosec */
	stmt := fmt.Sprintf("SELECT value FROM %s WHERE key = ?", table)
	results, err := session.Query(stmt, meta.key).WithContext(ctx).Iter().SliceMap()
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results for key %s", meta.key)
	}

	return types.NewResponse().
		SetData(results[0]["value"].([]byte)).
		SetMetadataKeyValue("key", meta.key), nil
}

func (c *Client) Exec(ctx context.Context, meta metadata, value []byte) (*types.Response, error) {
	session := c.session
	switch meta.consistency {
	case "One":
		sess, err := c.createSession(gocql.One)
		if err != nil {
			return nil, err
		}
		defer sess.Close()
		session = sess
	case "LocalOne":
		sess, err := c.createSession(gocql.LocalOne)
		if err != nil {
			return nil, err
		}
		defer sess.Close()
		session = sess
	case "LocalQuorum":
		sess, err := c.createSession(gocql.LocalQuorum)
		if err != nil {
			return nil, err
		}
		defer sess.Close()
		session = sess
	}
	query := string(value)
	if query == "" {
		return nil, fmt.Errorf("no query string found")
	}

	err := session.Query(query).WithContext(ctx).Exec()
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("result", "ok"),
		nil
}
func (c *Client) Query(ctx context.Context, meta metadata, value []byte) (*types.Response, error) {
	session := c.session
	switch meta.consistency {
	case "One":
		sess, err := c.createSession(gocql.One)
		if err != nil {
			return nil, err
		}
		defer sess.Close()
		session = sess
	case "LocalOne":
		sess, err := c.createSession(gocql.LocalOne)
		if err != nil {
			return nil, err
		}
		defer sess.Close()
		session = sess
	case "LocalQuorum":
		sess, err := c.createSession(gocql.LocalQuorum)
		if err != nil {
			return nil, err
		}
		defer sess.Close()
		session = sess
	}
	query := string(value)
	if query == "" {
		return nil, errors.New("no query string found")
	}
	results, err := session.Query(query).WithContext(ctx).Iter().SliceMap()
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("no results for this query")

	}
	return types.NewResponse().
		SetData(results[0]["value"].([]byte)).
		SetMetadataKeyValue("result", "ok"), nil

}
func (c *Client) Set(ctx context.Context, meta metadata, value []byte) (*types.Response, error) {
	session := c.session

	switch meta.consistency {
	case "One":
		sess, err := c.createSession(gocql.One)
		if err != nil {
			return nil, err
		}
		defer sess.Close()
		session = sess
	case "LocalOne":
		sess, err := c.createSession(gocql.LocalOne)
		if err != nil {
			return nil, err
		}
		defer sess.Close()
		session = sess
	case "LocalQuorum":
		sess, err := c.createSession(gocql.LocalQuorum)
		if err != nil {
			return nil, err
		}
		defer sess.Close()
		session = sess
	}

	table := meta.keyspaceTable()
	if table == "" {
		table = c.table
	}
	/* #nosec */
	stmt := fmt.Sprintf("INSERT INTO %s (key, value) VALUES (?, ?)", table)
	err := session.Query(stmt, meta.key, value).WithContext(ctx).Exec()
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("key", meta.key).
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) Delete(ctx context.Context, meta metadata) (*types.Response, error) {
	table := meta.keyspaceTable()
	if table == "" {
		table = c.table
	}
	/* #nosec */
	stmt := fmt.Sprintf("DELETE FROM %s WHERE key = ?", table)
	err := c.session.Query(stmt, meta.key).WithContext(ctx).Exec()
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("key", meta.key).
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) downloadFile() error {
	// Create the file
	out, err := os.Create("tls.pem")
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(c.opts.tls)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Stop() error {
	if c.session != nil {
		c.session.Close()
	}
	return nil
}
