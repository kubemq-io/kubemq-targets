package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"github.com/kubemq-hub/kubemq-target-connectors/pkg/logger"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
	"strconv"
	"strings"
	"time"

	redisClient "github.com/go-redis/redis/v7"
)

const (
	setQuery                 = "local var1 = redis.pcall(\"HGET\", KEYS[1], \"version\"); if type(var1) == \"table\" then redis.call(\"DEL\", KEYS[1]); end; if not var1 or type(var1)==\"table\" or var1 == \"\" or var1 == ARGV[1] or ARGV[1] == \"0\" then redis.call(\"HSET\", KEYS[1], \"data\", ARGV[2]) return redis.call(\"HINCRBY\", KEYS[1], \"version\", 1) else return error(\"failed to set key \" .. KEYS[1]) end"
	delQuery                 = "local var1 = redis.pcall(\"HGET\", KEYS[1], \"version\"); if not var1 or type(var1)==\"table\" or var1 == ARGV[1] or var1 == \"\" or ARGV[1] == \"0\" then return redis.call(\"DEL\", KEYS[1]) else return error(\"failed to delete \" .. KEYS[1]) end"
	connectedSlavesReplicas  = "connected_slaves:"
	infoReplicationDelimiter = "\r\n"
	defaultDB                = 0
)

// Client is a Client state store
type Client struct {
	name     string
	redis    *redisClient.Client
	opts     options
	replicas int
	log      *logger.Logger
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
	redisOpts := &redisClient.Options{
		Addr:            c.opts.host,
		Password:        c.opts.password,
		DB:              defaultDB,
		MaxRetries:      c.opts.maxRetries,
		MaxRetryBackoff: time.Duration(c.opts.maxRetryBackoffSeconds) * time.Second,
	}

	/* #nosec */
	if c.opts.enableTLS {
		redisOpts.TLSConfig = &tls.Config{
			InsecureSkipVerify: c.opts.enableTLS,
		}
	}

	c.redis = redisClient.NewClient(redisOpts)
	_, err = c.redis.WithContext(ctx).Ping().Result()
	if err != nil {
		return fmt.Errorf("error connecting to redis at %s: %w", c.opts.host, err)
	}
	c.replicas, err = c.getConnectedSlaves(ctx)
	return err
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

	}
	return nil, nil
}
func (c *Client) getConnectedSlaves(ctx context.Context) (int, error) {
	res, err := c.redis.DoContext(ctx, "INFO", "replication").Result()
	if err != nil {
		return 0, err
	}

	// Response example: https://redis.io/commands/info#return-value
	// # Replication\c\nrole:master\c\nconnected_slaves:1\c\n
	s, _ := strconv.Unquote(fmt.Sprintf("%q", res))
	if len(s) == 0 {
		return 0, nil
	}

	return c.parseConnectedSlaves(s), nil
}

func (c *Client) parseConnectedSlaves(res string) int {
	infos := strings.Split(res, infoReplicationDelimiter)
	for _, info := range infos {
		if strings.Contains(info, connectedSlavesReplicas) {
			parsedReplicas, _ := strconv.ParseUint(info[len(connectedSlavesReplicas):], 10, 32)
			return int(parsedReplicas)
		}
	}

	return 0
}
func (c *Client) Get(ctx context.Context, meta metadata) (*types.Response, error) {
	res, err := c.redis.DoContext(ctx, "HGETALL", meta.key).Result() // Prefer values with ETags
	if err != nil {
		return c.directGet(meta.key) //Falls back to original get
	}
	if res == nil {
		return types.NewResponse().
			SetMetadataKeyValue("key", meta.key).
			SetMetadataKeyValue("error", "true").
			SetMetadataKeyValue("message", "no data found for this key"), nil
	}
	vals := res.([]interface{})
	if len(vals) == 0 {
		return types.NewResponse().
			SetMetadataKeyValue("key", meta.key).
			SetMetadataKeyValue("error", "true").
			SetMetadataKeyValue("message", "no data found for this key"), nil
	}

	data, _, err := c.getKeyVersion(vals)
	if err != nil {
		return types.NewResponse().
			SetMetadataKeyValue("key", meta.key).
			SetMetadataKeyValue("error", "true").
			SetMetadataKeyValue("message", err.Error()), nil
	}
	return types.NewResponse().
		SetData([]byte(data)).
		SetMetadataKeyValue("error", "false").
		SetMetadataKeyValue("key", meta.key), nil
}

func (c *Client) getKeyVersion(vals []interface{}) (data string, version string, err error) {
	seenData := false
	seenVersion := false
	for i := 0; i < len(vals); i += 2 {
		field, _ := strconv.Unquote(fmt.Sprintf("%q", vals[i]))
		switch field {
		case "data":
			data, _ = strconv.Unquote(fmt.Sprintf("%q", vals[i+1]))
			seenData = true
		case "version":
			version, _ = strconv.Unquote(fmt.Sprintf("%q", vals[i+1]))
			seenVersion = true
		}
	}
	if !seenData || !seenVersion {
		return "", "", fmt.Errorf("required hash field 'data' or 'version' was not found")
	}
	return data, version, nil
}
func (c *Client) directGet(key string) (*types.Response, error) {
	res, err := c.redis.DoContext(context.Background(), "GET", key).Result()
	if err != nil {
		return nil, err
	}

	return types.NewResponse().
		SetMetadataKeyValue("key", key).
		SetMetadataKeyValue("error", "true").
		SetMetadataKeyValue("message", "no data found for this key"), nil

	s, _ := strconv.Unquote(fmt.Sprintf("%q", res))
	return types.NewResponse().
		SetMetadataKeyValue("key", key).
		SetMetadataKeyValue("error", "false").
		SetData([]byte(s)), nil
}

func (c *Client) Set(ctx context.Context, meta metadata, value []byte) (*types.Response, error) {

	if meta.concurrency == "last-write" {
		meta.etag = 0
	}

	_, err := c.redis.DoContext(ctx, "EVAL", setQuery, 1, meta.key, meta.etag, value).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to set key %s: %s", meta.key, err)
	}

	if meta.consistency == "strong" && c.replicas > 0 {
		_, err = c.redis.DoContext(ctx, "WAIT", c.replicas, 1000).Result()
		if err != nil {
			return nil, fmt.Errorf("timed out while waiting for %v replicas to acknowledge write", c.replicas)
		}
	}

	return types.NewResponse().
			SetMetadataKeyValue("key", meta.key).
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) Delete(ctx context.Context, meta metadata) (*types.Response, error) {
	_, err := c.redis.DoContext(ctx, "EVAL", delQuery, 1, meta.key, meta.etag).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to delete key '%s',%w", meta.key, err)
	}
	return types.NewResponse().
			SetMetadataKeyValue("key", meta.key).
			SetMetadataKeyValue("result", "ok"),
		nil
}
