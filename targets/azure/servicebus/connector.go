package servicebus

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.azure.servicebus").
		SetDescription("Azure Service Bus Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("end_point").
				SetDescription("Set Service Bus end point").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("shared_access_key_name").
				SetDescription("Set Service Bus shared access key name").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("shared_access_key").
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
		)
}
