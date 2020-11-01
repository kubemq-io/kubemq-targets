package mysql

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("gcp.stores.mysql").
		SetDescription("GCP MySQL Direct Mode Target").
		AddProperty(
			common.NewProperty().
				SetKind("condition").
				SetName("connection-type").
				SetDescription("Set MySQL Connection Type").
				SetMust(true).
				SetOptions([]string{"Proxy", "Direct"}).
				SetDefault("Proxy").
				NewCondition("Proxy", []*common.Property{
					common.NewProperty().
						SetKind("null").
						SetName("use_proxy").
						SetDescription("Set use proxy").
						SetMust(true).
						SetDefault("true"),
					common.NewProperty().
						SetKind("string").
						SetName("instance_connection_name").
						SetDescription("Set MySQL instance connection name").
						SetMust(true).
						SetDefault(""),
					common.NewProperty().
						SetKind("string").
						SetName("db_user").
						SetDescription("Set MySQL db user").
						SetMust(true).
						SetDefault(""),
					common.NewProperty().
						SetKind("string").
						SetName("db_password").
						SetDescription("Set MySQL db password").
						SetMust(true).
						SetDefault(""),
					common.NewProperty().
						SetKind("string").
						SetName("db_name").
						SetDescription("Sets Mysql db name").
						SetMust(true).
						SetDefault(""),
					common.NewProperty().
						SetKind("multilines").
				SetName("credentials").
						SetDescription("Set MySQL credentials").
						SetMust(true).
						SetDefault(""),
				}).
				NewCondition("Direct", []*common.Property{
					common.NewProperty().
						SetKind("null").
						SetName("use_proxy").
						SetDescription("Set use proxy").
						SetMust(true).
						SetDefault("false"),
					common.NewProperty().
						SetKind("string").
						SetName("connection").
						SetDescription("Set MySQL connection string").
						SetMust(true).
						SetDefault("root:mysql@(localhost:3306)/store?charset=utf8&parseTime=True&loc=Local"),
				}),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_idle_connections").
				SetDescription("Set MySQL max idle connections").
				SetMust(false).
				SetDefault("10").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_open_connections").
				SetDescription("Set MySQL max open connections").
				SetMust(false).
				SetDefault("100").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("connection_max_lifetime_seconds").
				SetDescription("Set MySQL connection max lifetime seconds").
				SetMust(false).
				SetDefault("3600").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set MySql execution method").
				SetOptions([]string{"query", "exec", "transaction"}).
				SetDefault("query").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("isolation_level").
				SetKind("string").
				SetDescription("Set MySql isolation level").
				SetOptions([]string{"Default", "ReadUncommitted", "ReadCommitted", "RepeatableRead", "Serializable"}).
				SetDefault("Default").
				SetMust(false),
		)
}
