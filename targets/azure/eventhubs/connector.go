package eventhubs

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("azure.eventhubs").
		SetDescription("Azure EventHubs Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("end_point").
				SetDescription("Set EventHubs end point").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("shared_access_key_name").
				SetDescription("Set EventHubs shared access key name").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("shared_access_key").
				SetDescription("Set EventHubs shared access key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("entity_path").
				SetDescription("Set EventHubs entity path").
				SetMust(true).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set EventHubs execution method").
				SetOptions([]string{"send", "send_batch"}).
				SetDefault("send").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("properties").
				SetKind("string").
				SetDescription("Set EventHubs properties").
				SetDefault("").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("partition_key").
				SetKind("string").
				SetDescription("Set EventHubs partition key").
				SetDefault("").
				SetMust(false),
		)
}
