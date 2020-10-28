package azuresql

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("azure.stores.azuresql").
		SetDescription("Azure SQL Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("connection").
				SetDescription("Set MSSQL connection string").
				SetMust(true).
				SetDefault("server=server.net;user id=test;password=test;port=1433;database=initial_db;"),
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
		)
}
