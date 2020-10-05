package rabbitmq

import (
	"github.com/kubemq-hub/builder/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.messaging.rabbitmq").
		SetDescription("RabbitMQ Messaging Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("gateway").
				SetDescription("Sets RabbitMQ url connection string").
				SetMust(true).
				SetDefault("amqp://rabbitmq:rabbitmq@localhost:5672/"),
		)
}
