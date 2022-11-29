package mysql

import (
	"math"

	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.rds.mysql").
		SetDescription("AWS RDS MySQL Target").
		SetName("MySQL").
		SetProvider("AWS").
		SetCategory("Store").
		SetTags("rds", "sql", "db", "cloud", "managed").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_key").
				SetDescription("Set MySQL aws key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_secret_key").
				SetDescription("Set MySQL aws secret key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("region").
				SetDescription("Set MySQL aws region").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set MySQL aws token").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("end_point").
				SetTitle("Endpoint").
				SetDescription("Set MySQL end point address").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("db_port").
				SetTitle("Port").
				SetDescription("Set MySQL end point port").
				SetMust(true).
				SetDefault("3306").
				SetMin(0).
				SetMax(65535),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("db_user").
				SetTitle("Username").
				SetDescription("Set MySQL db user").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("db_name").
				SetTitle("Database").
				SetDescription("Set MySQL db name").
				SetMust(true).
				SetDefault(""),
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
				SetTitle("Connection Lifetime (Seconds)").
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
