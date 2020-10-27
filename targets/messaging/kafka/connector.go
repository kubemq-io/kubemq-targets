package kafka

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("messaging.kafka").
		SetDescription("Kafka Messaging Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("brokers").
				SetDescription("Set Kafka brokers list").
				SetMust(true).
				SetDefault("localhost:9092"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sasl_username").
				SetDescription("Set Kafka username").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sasl_password").
				SetDescription("Set Kafka password").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("topic").
				SetDescription("Set Kafka topic").
				SetMust(true).
				SetDefault(""),
		)

}
