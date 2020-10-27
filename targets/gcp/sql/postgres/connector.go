package postgres

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.gcp.stores.postgres").
		SetDescription("GCP Postgres Direct Mode Target").
		AddProperty(
			common.NewProperty().
				SetKind("condition").
				SetName("connection-type").
				SetDescription("Sets Postgres Connection Type").
				SetMust(true).
				SetOptions([]string{"Proxy", "Direct"}).
				SetDefault("Proxy").
				NewCondition("Proxy", []*common.Property{
					common.NewProperty().
						SetKind("null").
						SetName("use_proxy").
						SetDescription("Sets use proxy").
						SetMust(true).
						SetDefault("true"),
					common.NewProperty().
						SetKind("string").
						SetName("instance_connection_name").
						SetDescription("Sets Postgres instance connection name").
						SetMust(true).
						SetDefault(""),
					common.NewProperty().
						SetKind("string").
						SetName("db_user").
						SetDescription("Sets Postgres db user").
						SetMust(true).
						SetDefault(""),
					common.NewProperty().
						SetKind("string").
						SetName("db_password").
						SetDescription("Sets Postgres db password").
						SetMust(true).
						SetDefault(""),
					common.NewProperty().
						SetKind("string").
						SetName("db_name").
						SetDescription("Sets Postgres db name").
						SetMust(true).
						SetDefault(""),
					common.NewProperty().
						SetKind("string").
						SetName("credentials").
						SetDescription("Sets GCP credentials").
						SetMust(true).
						SetDefault(""),
				}).
				NewCondition("Direct", []*common.Property{
					common.NewProperty().
						SetKind("null").
						SetName("use_proxy").
						SetDescription("Sets use proxy").
						SetMust(true).
						SetDefault("false"),
					common.NewProperty().
						SetKind("string").
						SetName("connection").
						SetDescription("Sets Postgres connection string").
						SetMust(true).
						SetDefault("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
				}),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_idle_connections").
				SetDescription("Sets Postgres max idle connections").
				SetMust(false).
				SetDefault("10").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_open_connections").
				SetDescription("Sets Postgres max open connections").
				SetMust(false).
				SetDefault("100").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("connection_max_lifetime_seconds").
				SetDescription("Sets Postgres connection max lifetime seconds").
				SetMust(false).
				SetDefault("3600").
				SetMin(1).
				SetMax(math.MaxInt32),
		)
}
