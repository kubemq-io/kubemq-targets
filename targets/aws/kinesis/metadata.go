package kinesis

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

const (
	DefaultLimit      = 1
	DefaultShardCount = 1
)

type metadata struct {
	method string
	
	shardCount        int64
	streamName        string
	partitionKey      string
	shardID           string
	limit             int64
	streamARN         string
	shardIteratorType string
	shardIteratorID string
}

var methodsMap = map[string]string{
	"list_streams":          "list_streams",
	"list_stream_consumers": "list_stream_consumers",
	"create_stream":         "create_stream",
	"delete_stream":         "delete_stream",
	"put_record":            "put_record",
	"put_records":           "put_records",
	"get_records":            "get_records",
	"get_shard_iterator":    "get_shard_iterator",
	"list_shards":           "list_shards",
}


func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidHttpMethodTypes(methodsMap)
	}
	m.shardCount = int64(meta.ParseInt("shard_count", DefaultShardCount))
	m.limit = int64(meta.ParseInt("limit", DefaultLimit))
	if m.method != "list_streams" && m.method != "list_stream_consumers" && m.method != "get_record" {
		m.streamName, err = meta.MustParseString("stream_name")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing stream_name, %w", err)
		}
	}
	switch m.method {
	case "put_record":
		m.partitionKey,err = meta.MustParseString("partition_key")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing partition_key, %w", err)
		}
	case "get_shard_iterator":
		m.shardID,err = meta.MustParseString("shard_id")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing shard_id, %w", err)
		}
		m.shardIteratorType ,err = meta.MustParseString("shard_iterator_type")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing shard_iterator_type, %w", err)
		}
	case "get_records":
		m.shardIteratorID,err = meta.MustParseString("shard_iterator_id")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing shard_iterator_id, %w", err)
		}
	case "list_stream_consumers":
		m.streamARN,err = meta.MustParseString("stream_arn")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing stream_arn, %w", err)
		}
	}

	return m, nil
}
