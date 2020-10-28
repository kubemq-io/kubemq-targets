package blob

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("azure.storage.blob").
		SetDescription("Azure Blob Storage Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("storage_access_key").
				SetDescription("Set Blob Storage storage access key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("storage_account").
				SetDescription("Set Blob Storage storage account").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("policy").
				SetDescription("Set Blob Storage retry policy").
				SetOptions([]string{"exponential", "fixed"}).
				SetMust(true).
				SetDefault("exponential"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_tries").
				SetDescription("Set Blob Storage max tries").
				SetMust(false).
				SetDefault("1").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("try_timeout").
				SetDescription("Set Blob Storage try timeout in milliseconds").
				SetMust(false).
				SetDefault("1000").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("retry_delay").
				SetDescription("Set Blob Storage retry delay in milliseconds").
				SetMust(false).
				SetDefault("60000").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_retry_delay").
				SetDescription("Set Blob Storage max retry delay in milliseconds").
				SetMust(false).
				SetDefault("180000").
				SetMin(1).
				SetMax(math.MaxInt32),
		)
}
