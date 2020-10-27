package cassandra

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("stores.cassandra").
		SetDescription("Cassandra Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("host").
				SetDescription("Set Cassandra hosts addresses").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("port").
				SetDescription("Set Cassandra port").
				SetMust(true).
				SetMin(0).
				SetMax(65535).
				SetDefault("9042"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Set Cassandra username").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Set Cassandra password").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("replication_factor").
				SetDescription("Set Cassandra replication factor").
				SetMust(false).
				SetDefault("1").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("consistency").
				SetDescription("Set Cassandra consistency").
				SetMust(true).
				SetOptions([]string{"All", "One", "Two", "Three", "Quorum", "LocalQuorum", "EachQuorum", "LocalOne", "Any"}).
				SetDefault("All"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("default_table").
				SetDescription("Set Cassandra default table").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("default_keyspace").
				SetDescription("Set Cassandra default keyspace").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("connect_timeout_seconds").
				SetDescription("Set Cassandra connection timeout in seconds").
				SetMust(false).
				SetDefault("60").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("timeout_seconds").
				SetDescription("Set Cassandra operation timeout in seconds").
				SetMust(false).
				SetDefault("60").
				SetMin(1).
				SetMax(math.MaxInt32),
		)
}
