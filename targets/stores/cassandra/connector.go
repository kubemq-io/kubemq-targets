package cassandra

import (
	"math"

	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("stores.cassandra").
		SetDescription("Cassandra Target").
		SetName("Cassandra").
		SetProvider("").
		SetCategory("Store").
		SetTags("db", "sql", "no-sql").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("hosts").
				SetTitle("Hosts Addresses").
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
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Set Cassandra password").
				SetMust(false).
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
				SetMust(false).
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
				SetTitle("Connect Timeout (Seconds)").
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
				SetTitle("Operation Timeout (Seconds)").
				SetDescription("Set Cassandra operation timeout in seconds").
				SetMust(false).
				SetDefault("60").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Cassandra execution method").
				SetOptions([]string{"get", "set", "delete", "query", "exec"}).
				SetDefault("get").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("consistency").
				SetDescription("Set Cassandra consistency Level").
				SetOptions([]string{"strong", "eventual", ""}).
				SetDefault("strong").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("key").
				SetKind("string").
				SetDescription("Cassandra key to set get or delete").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("table").
				SetKind("string").
				SetDescription("Cassandra table name").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("keyspace").
				SetKind("string").
				SetDescription("Cassandra keyspace data container name").
				SetDefault("").
				SetMust(false),
		)
}
