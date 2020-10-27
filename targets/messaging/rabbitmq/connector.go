package rabbitmq

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("messaging.rabbitmq").
		SetDescription("RabbitMQ Messaging Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("gateway").
				SetDescription("Set RabbitMQ url connection string").
				SetMust(true).
				SetDefault("amqp://rabbitmq:rabbitmq@localhost:5672/"),
		)
}
