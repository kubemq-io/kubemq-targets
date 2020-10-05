package redis

import (
	"context"
	"fmt"
	redisClient "github.com/go-redis/redis/v7"
	"github.com/kubemq-hub/builder/common"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/types"
	"strconv"
	"strings"
)

const (
	setQuery                 = "local var1 = redis.pcall(\"HGET\", KEYS[1], \"version\"); if type(var1) == \"table\" then redis.call(\"DEL\", KEYS[1]); end; if not var1 or type(var1)==\"table\" or var1 == \"\" or var1 == ARGV[1] or ARGV[1] == \"0\" then redis.call(\"HSET\", KEYS[1], \"data\", ARGV[2]) return redis.call(\"HINCRBY\", KEYS[1], \"version\", 1) else return error(\"failed to set key \" .. KEYS[1]) end"
	delQuery                 = "local var1 = redis.pcall(\"HGET\", KEYS[1], \"version\"); if not var1 or type(var1)==\"table\" or var1 == ARGV[1] or var1 == \"\" or ARGV[1] == \"0\" then return redis.call(\"DEL\", KEYS[1]) else return error(\"failed to delete \" .. KEYS[1]) end"
	connectedSlavesReplicas  = "connected_slaves:"
	infoReplicationDelimiter = "\r\n"
)

// Client is a Client state store
type Client struct {
	name     string
	redis    *redisClient.Client
	opts     options
	replicas int
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
	redisInfo, err := redisClient.ParseURL(c.opts.url)
	if err != nil {
		return fmt.Errorf("error parsing redis url %s: %w", c.opts.url, err)
	}
	c.redis = redisClient.NewClient(redisInfo)
	_, err = c.redis.WithContext(ctx).Ping().Result()
	if err != nil {
		return fmt.Errorf("error connecting to redis at %s: %w", redisInfo.Addr, err)
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
		return c.directGet(ctx, meta.key) //Falls back to original get
	}
	if res == nil {
		return nil, fmt.Errorf("no data found for this key")
	}
	vals := res.([]interface{})
	if len(vals) == 0 {
		return nil, fmt.Errorf("no data found for this key")
	}

	data, _, err := c.getKeyVersion(vals)
	if err != nil {
		return nil, fmt.Errorf("error found for get this key, %w", err)
	}
	return types.NewResponse().
		SetData([]byte(data)).
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
func (c *Client) directGet(ctx context.Context, key string) (*types.Response, error) {
	res, err := c.redis.DoContext(ctx, "GET", key).Result()
	if err != nil {
		return nil, err
	}
	s, _ := strconv.Unquote(fmt.Sprintf("%q", res))
	return types.NewResponse().
		SetMetadataKeyValue("key", key).
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
