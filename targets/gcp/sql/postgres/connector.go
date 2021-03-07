package postgres

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("gcp.stores.postgres").
		SetDescription("GCP Postgres Direct Mode Target").
		SetName("Postgres").
		SetProvider("GCP").
		SetCategory("Store").
		SetTags("db","sql","cloud","managed").
		AddProperty(
			common.NewProperty().
				SetKind("condition").
				SetName("connection-type").
				SetTitle("Connection Type").
				SetDescription("Set Postgres Connection Type").
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
						SetDescription("Set Postgres instance connection name").
						SetMust(true).
						SetDefault(""),
					common.NewProperty().
						SetKind("string").
						SetName("db_user").
						SetTitle("Username").
						SetDescription("Set Postgres db user").
						SetMust(true).
						SetDefault(""),
					common.NewProperty().
						SetKind("string").
						SetName("db_password").
						SetTitle("Password").
						SetDescription("Set Postgres db password").
						SetMust(true).
						SetDefault(""),
					common.NewProperty().
						SetKind("multilines").
				SetName("credentials").
						SetDescription("Set Postgres credentials").
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
						SetTitle("Connection String").
						SetDescription("Set Postgres connection string").
						SetMust(true).
						SetDefault("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
				}),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_idle_connections").
				SetDescription("Set Postgres max idle connections").
				SetMust(false).
				SetDefault("10").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_open_connections").
				SetDescription("Set Postgres max open connections").
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
				SetDescription("Set Postgres connection max lifetime seconds").
				SetMust(false).
				SetDefault("3600").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Postgres execution method").
				SetOptions([]string{"query", "exec", "transaction"}).
				SetDefault("query").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("isolation_level").
				SetKind("string").
				SetDescription("Set Postgres isolation level").
				SetOptions([]string{"Default", "ReadUncommitted", "ReadCommitted", "RepeatableRead", "Serializable"}).
				SetDefault("Default").
				SetMust(false),
		)
}
