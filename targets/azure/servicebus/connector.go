package servicebus

import (
	"math"

	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("azure.servicebus").
		SetDescription("Azure Service Bus Target").
		SetName("ServiceBus").
		SetProvider("Azure").
		SetCategory("Messaging").
		SetTags("queue", "pub/sub", "cloud", "managed").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("end_point").
				SetTitle("Endpoint").
				SetDescription("Set Service Bus end point").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("shared_access_key_name").
				SetTitle("Access Key Name").
				SetDescription("Set Service Bus shared access key name").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("shared_access_key").
				SetTitle("Access Key").
				SetDescription("Set Service Bus shared access key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("queue_name").
				SetDescription("Set Service Bus queue name").
				SetMust(true).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Service Bus execution method").
				SetOptions([]string{"send", "send_batch"}).
				SetDefault("send").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("label").
				SetKind("string").
				SetDescription("Set Service Bus label").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("content_type").
				SetKind("string").
				SetDescription("Set Service Bus content type").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("time_to_live").
				SetDescription("Set Blob Storage time to live milliseconds").
				SetMust(false).
				SetDefault("1000000000").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("max_batch_size").
				SetDescription("Set Blob Storage max batch size in bytes").
				SetMust(false).
				SetDefault("1024").
				SetMin(0).
				SetMax(math.MaxInt32),
		)
}
