package bigquery

import "github.com/kubemq-hub/builder/connector/common"

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("gcp.bigquery").
		SetDescription("GCP Bigquery Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("project_id").
				SetDescription("Set GCP project ID").
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
				SetDescription("Set GCP BigQuery execution method").
				SetOptions([]string{"query", "create_data_set", "delete_data_set", "create_table","delete_table", "get_table_info", "get_data_sets", "insert"}).
				SetDefault("query").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("table_name").
				SetKind("string").
				SetDescription("Set BigQuery table name").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("dataset_id").
				SetKind("string").
				SetDescription("Set BigQuery dataset id").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("location").
				SetKind("string").
				SetDescription("Set BigQuery location").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("query").
				SetKind("string").
				SetDescription("Set BigQuery query").
				SetDefault("").
				SetMust(false),
		)
}
