package eventhubs

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.azure.eventhubs").
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
		)
}
