package keyspaces

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.aws.keyspaces").
		SetDescription("AWS Keyspaces Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("hosts").
				SetDescription("Set Keyspaces hosts addresses").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("port").
				SetDescription("Set Keyspaces port").
				SetMust(true).
				SetMin(0).
				SetMax(65535).
				SetDefault("9142"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Set Keyspaces username").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Set Keyspaces password").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("tls").
				SetDescription("Set Keyspaces tls download url").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("replication_factor").
				SetDescription("Set Keyspaces replication factor").
				SetMust(false).
				SetDefault("1").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("proto_version").
				SetDescription("Set Keyspaces protoVersion").
				SetMust(false).
				SetDefault("4").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("consistency").
				SetDescription("Set Keyspaces consistency").
				SetMust(true).
				SetOptions([]string{"One", "LocalQuorum", "local_one"}).
				SetDefault("LocalQuorum"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("default_table").
				SetDescription("Set Keyspaces default table").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("default_keyspace").
				SetDescription("Set Keyspaces default keyspace").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("connect_timeout_seconds").
				SetDescription("Set Keyspaces connection timeout in seconds").
				SetMust(false).
				SetDefault("60").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("timeout_seconds").
				SetDescription("Set Keyspaces operation timeout in seconds").
				SetMust(false).
				SetDefault("60").
				SetMin(1).
				SetMax(math.MaxInt32),
		)
}
