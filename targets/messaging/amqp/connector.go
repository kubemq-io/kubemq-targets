package amqp

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("messaging.amqp").
		SetDescription("AMQP Messaging Target").
		SetName("AMQP").
		SetProvider("").
		SetCategory("Messaging").
		SetTags("address", "pub/sub", "queues").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("url").
				SetDescription("Set AMQP url connection string").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Set AMQP username").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Set AMQP password").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("skip_insecure").
				SetDescription("(SSL) Set skip TLS Certificate verification").
				SetMust(false).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("ca_cert").
				SetDescription("(SSL) Set TLS CA Certificate").
				SetMust(false).
				SetDefault(""),
		)
}
