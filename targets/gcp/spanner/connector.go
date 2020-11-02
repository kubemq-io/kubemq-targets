package spanner

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("gcp.spanner").
		SetDescription("GCP Spanner Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("db").
				SetDescription("Set GCP Spanner DB").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("credentials").
				SetDescription("Set GCP credentials").
				SetMust(true).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Spanner execution method").
				SetOptions([]string{"query", "read", "update_database_ddl", "insert", "update", "insert_or_update"}).
				SetDefault("query").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("query").
				SetKind("string").
				SetDescription("Set Spanner query request").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("table_name").
				SetKind("string").
				SetDescription("Set Spanner table_name").
				SetDefault("").
				SetMust(false),
		)
}
