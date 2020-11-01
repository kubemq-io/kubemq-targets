package bigtable

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("gcp.bigtable").
		SetDescription("GCP Bigtable Target").
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
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("instance").
				SetDescription("Set Bigtable instance").
				SetMust(true).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set GCP Bigtable execution method").
				SetOptions([]string{"write", "write_batch", "get_row", "get_all_rows","delete_row", "get_tables", "create_table", "delete_table", "create_column_family", "get_all_rows_by_column"}).
				SetDefault("write").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("table_name").
				SetKind("string").
				SetDescription("Set Bigtable table name").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("column_family").
				SetKind("string").
				SetDescription("Set Bigtable column family").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("row_key_prefix").
				SetKind("string").
				SetDescription("Set Bigtable row key prefix").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("column_name").
				SetKind("string").
				SetDescription("Set Bigtable column name").
				SetDefault("").
				SetMust(false),
		)
}
