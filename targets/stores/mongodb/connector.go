package mongodb

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("stores.mongodb").
		SetDescription("MongoDB Target").
		SetName("MongoDB").
		SetProvider("").
		SetCategory("Store").
		SetTags("db", "no-sql").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("host").
				SetTitle("Host address").
				SetDescription("Set MongoDB host address").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Set MongoDB username").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Set MongoDB password").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("database").
				SetDescription("Set MongoDB database").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("collection").
				SetDescription("Set MongoDB collection").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("params").
				SetDescription("Set MongoDB params").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("read_concurrency").
				SetDescription("Set MongoDB read concurrency").
				SetOptions([]string{"local", "majority", "available", "linearizable", "snapshot"}).
				SetMust(false).
				SetDefault("local"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("write_concurrency").
				SetDescription("Set MongoDB write concurrency").
				SetOptions([]string{"majority", "Other"}).
				SetMust(false).
				SetDefault("majority"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("operation_timeout_seconds").
				SetTitle("Operation Timeout (Seconds)").
				SetDescription("Set MongoDB operation timeout seconds").
				SetMust(false).
				SetDefault("90").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set MongoDB execution method").
				SetOptions([]string{"get_by_key", "set_by_key", "delete_by_key", "find", "find_many",
					"insert", "insert_many", "update", "update_many", "delete_one",
					"delete_many", "aggregate", "distinct",
				}).
				SetDefault("get").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("key").
				SetKind("string").
				SetDescription("Set MongoDB key").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("filter").
				SetKind("string").
				SetDescription("Set filter").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("set_upsert").
				SetKind("bool").
				SetDescription("Set Upsert in update mode").
				SetMust(false),
		)

}
