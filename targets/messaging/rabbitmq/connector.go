package rabbitmq

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("messaging.rabbitmq").
		SetDescription("RabbitMQ Messaging Target").
		SetName("RabbitMQ").
		SetProvider("").
		SetCategory("Messaging").
		SetTags("queue", "pub/sub").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("url").
				SetDescription("Set RabbitMQ url connection string").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("insecure").
				SetDescription("Set skip TLS Certificate verification").
				SetMust(false).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("ca_cert").
				SetDescription("Set TLS CA Certificate").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("client_certificate").
				SetDescription("Set TLS Client PEM data").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("client_key").
				SetDescription("Set TLS Client Key PEM data").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("default_exchange").
				SetDescription("Set Default Exchange for routing").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("default_topic").
				SetDescription("Set Default Topic for routing").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("default_persistence").
				SetDescription("Set Default Persistence for routed message").
				SetMust(false).
				SetDefault("true"),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("queue").
				SetKind("string").
				SetDescription("Set RabbitMQ queue Name").
				SetDefault("").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("exchange").
				SetKind("string").
				SetDescription("Set RabbitMQ exchange name").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("bool").
				SetName("mandatory").
				SetDescription("Set RabbitMQ mandatory").
				SetMust(false).
				SetDefault("false"),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("bool").
				SetName("immediate").
				SetDescription("Set RabbitMQ immediate").
				SetMust(false).
				SetDefault("false"),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("delivery_mode").
				SetDescription("Set RabbitMQ delivery mode").
				SetMust(true).
				SetMin(0).
				SetMax(2).
				SetDefault("1"),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("priority").
				SetDescription("Set RabbitMQ priority").
				SetMust(true).
				SetMin(0).
				SetMax(9).
				SetDefault("0"),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("correlation_id").
				SetKind("string").
				SetDescription("Set RabbitMQ correlation id ").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("reply_to").
				SetKind("string").
				SetDescription("Set RabbitMQ set reply to target ").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("expiry_seconds").
				SetDescription("Set RabbitMQ expiry in seconds").
				SetMust(true).
				SetMin(0).
				SetMax(math.MaxInt32),
		)
}
