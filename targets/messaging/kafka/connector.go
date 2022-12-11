package kafka

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("messaging.kafka").
		SetDescription("Kafka Messaging Target").
		SetName("Kafka").
		SetProvider("").
		SetCategory("Messaging").
		SetTags("streaming", "pub/sub").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("brokers").
				SetTitle("Brokers Address").
				SetDescription("Set Kafka brokers list").
				SetMust(true).
				SetDefault("localhost:9092"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sasl_username").
				SetTitle("SASL Username").
				SetDescription("Set Kafka username").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sasl_password").
				SetTitle("SASL Password").
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
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sasl_mechanism").
				SetDescription("Set SASL Mechanism").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("security_protocol").
				SetDescription("Set Security Protocol").
				SetMust(false).
				SetDefault(""),
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
				SetKind("bool").
				SetName("insecure").
				SetDescription("Set self-signed SSL Certificate").
				SetMust(false).
				SetDefault("false"),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("headers").
				SetKind("string").
				SetDescription("Set Kafka headers").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("key").
				SetKind("string").
				SetDescription("Set Kafka Key").
				SetDefault("").
				SetMust(false),
		)
}
