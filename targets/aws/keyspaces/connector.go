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
				SetName("host").
				SetDescription("Sets Keyspaces hosts addresses").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("tls").
				SetDescription("Sets Keyspaces tls download url").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("port").
				SetDescription("Sets Keyspaces port").
				SetMust(true).
				SetMin(0).
				SetMax(65535).
				SetDefault("9042"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Sets Keyspaces username").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Sets Keyspaces password").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("replication_factor").
				SetDescription("Sets Keyspaces replication factor").
				SetMust(false).
				SetDefault("1").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("consistency").
				SetDescription("Sets Keyspaces consistency").
				SetMust(true).
				SetOptions([]string{"All", "One", "Two", "Three", "Quorum", "LocalQuorum", "EachQuorum", "LocalOne", "Any"}).
				SetDefault("All"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("default_table").
				SetDescription("Sets Keyspaces default table").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("default_keyspace").
				SetDescription("Sets Keyspaces default keyspace").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("connect_timeout_seconds").
				SetDescription("Sets Keyspaces connection timeout in seconds").
				SetMust(false).
				SetDefault("60").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("timeout_seconds").
				SetDescription("Sets Keyspaces operation timeout in seconds").
				SetMust(false).
				SetDefault("60").
				SetMin(1).
				SetMax(math.MaxInt32),
		)
}
