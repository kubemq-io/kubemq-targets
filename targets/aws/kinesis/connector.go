package kinesis

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.kinesis").
		SetDescription("AWS Kinesis Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_key").
				SetDescription("Set Kinesis aws key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_secret_key").
				SetDescription("Set Kinesis aws secret key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("region").
				SetDescription("Set Kinesis aws region").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set Kinesis aws token").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Kinesis execution method").
				SetOptions([]string{"list_streams", "list_stream_consumers", "create_stream", "delete_stream", "put_record", "put_records", "get_records", "get_shard_iterator", "list_shards"}).
				SetDefault("put_record").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("stream_name").
				SetKind("string").
				SetDescription("Set Kinesis stream name").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("partition_key").
				SetKind("string").
				SetDescription("Set Kinesis partition key").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("shard_id").
				SetKind("string").
				SetDescription("Set Kinesis shard id").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("stream_arn").
				SetKind("string").
				SetDescription("Set Kinesis stream arn").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("shard_iterator_type").
				SetKind("string").
				SetDescription("Set Kinesis shard iterator type").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("shard_iterator_id").
				SetKind("string").
				SetDescription("Set Kinesis shard iterator id").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("limit").
				SetDescription("Set Kinesis limit").
				SetMust(false).
				SetDefault("1").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("shard_count").
				SetDescription("Set Kinesis shard count").
				SetMust(false).
				SetDefault("1").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("consumer_name").
				SetDescription("Set consumer name").
				SetDefault("").
				SetMust(false),
		)
}
