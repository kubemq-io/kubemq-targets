package cassandra

import (
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
)

// Client is a Client state store
type Client struct {
	name    string
	session *gocql.Session
	cluster *gocql.ClusterConfig
	table   string
	opts    options
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
	c.cluster = gocql.NewCluster(c.opts.hosts...)

	if c.opts.username != "" && c.opts.password != "" {
		c.cluster.Authenticator = gocql.PasswordAuthenticator{Username: c.opts.username, Password: c.opts.password}
	}
	c.cluster.Port = c.opts.port
	c.cluster.ProtoVersion = c.opts.protoVersion
	c.cluster.Consistency = c.opts.consistency
	c.cluster.Timeout = c.opts.timeoutSeconds
	c.cluster.ConnectTimeout = c.opts.connectTimeoutSeconds
	session, err := c.cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("error creating session: %s", err)
	}
	c.session = session
	if c.opts.defaultKeyspace != "" && c.opts.defaultTable != "" {
		err = c.tryCreateKeyspace(c.opts.defaultKeyspace, c.opts.replicationFactor)
		if err != nil {
			return fmt.Errorf("error creating defaultKeyspace %s: %s", c.opts.defaultTable, err)
		}
		err = c.tryCreateTable(c.opts.defaultTable, c.opts.defaultKeyspace)
		if err != nil {
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
	case "strong":
		sess, err := c.createSession(gocql.All)
		if err != nil {
			return nil, err
		}
		defer sess.Close()
		session = sess
	case "eventual":
		sess, err := c.createSession(gocql.One)
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
	case "strong":
		sess, err := c.createSession(gocql.Quorum)
		if err != nil {
			return nil, err
		}
		defer sess.Close()
		session = sess
	case "eventual":
		sess, err := c.createSession(gocql.Any)
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
	case "strong":
		sess, err := c.createSession(gocql.Quorum)
		if err != nil {
			return nil, err
		}
		defer sess.Close()
		session = sess
	case "eventual":
		sess, err := c.createSession(gocql.Any)
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
	results, err := session.Query(query).WithContext(ctx).Iter().SliceMap()
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no results for this query")

	}
	return types.NewResponse().
		SetData(results[0]["value"].([]byte)).
		SetMetadataKeyValue("result", "ok"), nil

}
func (c *Client) Set(ctx context.Context, meta metadata, value []byte) (*types.Response, error) {
	session := c.session

	switch meta.consistency {
	case "strong":
		sess, err := c.createSession(gocql.Quorum)
		if err != nil {
			return nil, err
		}
		defer sess.Close()
		session = sess
	case "eventual":
		sess, err := c.createSession(gocql.Any)
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
