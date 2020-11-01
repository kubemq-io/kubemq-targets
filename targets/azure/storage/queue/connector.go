package queue

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("azure.storage.queue").
		SetDescription("Azure Queue Storage Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("storage_access_key").
				SetDescription("Set Queue Storage storage access key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("storage_account").
				SetDescription("Set Queue Storage storage account").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("policy").
				SetDescription("Set Queue Storage retry policy").
				SetOptions([]string{"retry_policy_exponential", "retry_policy_fixed"}).
				SetMust(true).
				SetDefault("retry_policy_exponential"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_tries").
				SetDescription("Set Queue Storage max tries").
				SetMust(false).
				SetDefault("1").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("try_timeout").
				SetDescription("Set Queue Storage try timeout in milliseconds").
				SetMust(false).
				SetDefault("1000").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("retry_delay").
				SetDescription("Set Queue Storage retry delay in milliseconds").
				SetMust(false).
				SetDefault("60000").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_retry_delay").
				SetDescription("Set Queue Storage max retry delay in milliseconds").
				SetMust(false).
				SetDefault("180000").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Queue Storage execution method").
				SetOptions([]string{"query", "create_data_set", "delete_data_set", "create_table","delete_table", "get_table_info", "get_data_sets", "insert"}).
				SetDefault("query").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("query").
				SetKind("string").
				SetDescription("Set Queue Storage query").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("table_name").
				SetKind("string").
				SetDescription("Set Queue Storage table name").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("dataset_id").
				SetKind("string").
				SetDescription("Set Queue Storage dataset id").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("location").
				SetKind("string").
				SetDescription("Set Queue Storage location").
				SetDefault("").
				SetMust(false),
		)
}
