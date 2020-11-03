package mssql

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("stores.mssql").
		SetDescription("MSSQL Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("connection").
				SetDescription("Set MSSQL connection string").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_idle_connections").
				SetDescription("Set MSSQL max idle connections").
				SetMust(false).
				SetDefault("10").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_open_connections").
				SetDescription("Set MSSQL max open connections").
				SetMust(false).
				SetDefault("100").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("connection_max_lifetime_seconds").
				SetDescription("Set MSSQL connection max lifetime seconds").
				SetMust(false).
				SetDefault("3600").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set MSSQL execution method").
				SetOptions([]string{"query", "exec", "transaction"}).
				SetDefault("query").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("isolation_level").
				SetKind("string").
				SetDescription("Set MSSQL isolation level").
				SetOptions([]string{"Default", "ReadUncommitted", "ReadCommitted", "RepeatableRead", "Serializable"}).
				SetDefault("Default").
				SetMust(false),
		)

}
