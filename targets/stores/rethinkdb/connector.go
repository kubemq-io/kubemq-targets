package rethinkdb

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("stores.rethinkdb").
		SetDescription("Rethinkdb Target").
		SetName("RethinkDB").
		SetProvider("").
		SetCategory("Store").
		SetTags("db","sql").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("host").
				SetDescription("Set Rethinkdb host address").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Set Rethinkdb username").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Set Rethinkdb password").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("timeout").
				SetDescription("Set Rethinkdb operation timeout seconds").
				SetMust(false).
				SetDefault("5").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("keep_alive_period").
				SetDescription("Set Rethinkdb operation keep alive period seconds").
				SetMust(false).
				SetDefault("5").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("auth_key").
				SetDescription("Set Rethinkdb auth key").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("condition").
				SetName("tls").
				SetOptions([]string{"true", "false"}).
				SetDescription("Set tls conditions").
				SetMust(true).
				SetDefault("false").
				NewCondition("true", []*common.Property{
					common.NewProperty().
						SetKind("multilines").
						SetName("cert_key").
						SetDescription("Set certificate key").
						SetMust(false).
						SetDefault(""),
					common.NewProperty().
						SetKind("multilines").
						SetName("cert_file").
						SetDescription("Set certificate file").
						SetMust(false).
						SetDefault(""),
				}),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("hand_shake_version").
				SetDescription("Set Rethinkdb hand shake version").
				SetMust(false).
				SetDefault("0").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("number_of_retries").
				SetDescription("Set Rethinkdb number of retries").
				SetMust(false).
				SetDefault("0").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("initial_cap").
				SetDescription("Set Rethinkdb initial cap").
				SetMust(false).
				SetDefault("0").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_open").
				SetDescription("Set Rethinkdb max open connections").
				SetMust(false).
				SetDefault("0").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Rethinkdb execution method").
				SetOptions([]string{"get", "update", "delete", "insert"}).
				SetDefault("get").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("key").
				SetKind("string").
				SetDescription("Set Rethinkdb key").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("table").
				SetKind("string").
				SetDescription("Set Rethinkdb table name").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("db_name").
				SetKind("string").
				SetDescription("Set Rethinkdb db name").
				SetMust(true),
		)
}
