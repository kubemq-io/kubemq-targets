package queue

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("azure.storage.queue").
		SetDescription("Azure Queue Storage Target").
		SetName("Queue").
		SetProvider("Azure").
		SetCategory("Storage").
		SetTags("queue","messaging","db","cloud","managed").
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
				SetTitle("Try Timout (milliseconds)").
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
				SetOptions([]string{"create", "get_messages_count", "peek", "push","pop", "delete"}).
				SetDefault("create").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("queue_name").
				SetKind("string").
				SetDescription("Set Queue Storage queue name").
				SetDefault("").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("service_url").
				SetKind("string").
				SetDescription("Set Queue Storage service url").
				SetDefault("").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("queue_metadata").
				SetKind("string").
				SetDescription("Set Queue Storage queue metadata").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("max_messages").
				SetDescription("Set Queue Storage max messages").
				SetMust(false).
				SetDefault("32").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("visibility_timeout").
				SetDescription("Set Queue Storage visibility timeout").
				SetMust(false).
				SetDefault("100000000").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("time_to_live").
				SetDescription("Set Queue Storage time to live").
				SetMust(false).
				SetDefault("100000000").
				SetMin(1).
				SetMax(math.MaxInt32),
		)
}
