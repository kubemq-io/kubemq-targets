package athena

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.athena").
		SetDescription("AWS Athena Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_key").
				SetDescription("Set Athena aws key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_secret_key").
				SetDescription("Set Athena aws secret key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("region").
				SetDescription("Set Athena aws region").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set Athena aws token").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Athena execution method").
				SetOptions([]string{"list_databases", "list_data_catalogs", "query", "get_query_result"}).
				SetDefault("query").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("query").
				SetKind("string").
				SetDescription("Set Athena query").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("catalog").
				SetKind("string").
				SetDescription("Set Athena catalog").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("db").
				SetKind("string").
				SetDescription("Set Athena db").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("output_location").
				SetKind("string").
				SetDescription("Set Athena output location").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("execution_id").
				SetKind("string").
				SetDescription("Set Athena execution id").
				SetDefault("").
				SetMust(false),
		)
}
