package files

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("azure.storage.files").
		SetDescription("Azure Files Storage Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("storage_access_key").
				SetDescription("Set Files Storage storage access key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("storage_account").
				SetDescription("Set Files Storage storage account").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("policy").
				SetDescription("Set Files Storage retry policy").
				SetOptions([]string{"exponential", "fixed"}).
				SetMust(true).
				SetDefault("exponential"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_tries").
				SetDescription("Set Files Storage max tries").
				SetMust(false).
				SetDefault("1").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("try_timeout").
				SetDescription("Set Files Storage try timeout in milliseconds").
				SetMust(false).
				SetDefault("1000").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("retry_delay").
				SetDescription("Set Files Storage retry delay in milliseconds").
				SetMust(false).
				SetDefault("60000").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_retry_delay").
				SetDescription("Set Files Storage max retry delay in milliseconds").
				SetMust(false).
				SetDefault("180000").
				SetMin(1).
				SetMax(math.MaxInt32),
		)
}
