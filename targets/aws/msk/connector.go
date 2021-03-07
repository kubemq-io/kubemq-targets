package msk

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.msk").
		SetDescription("AWS MSK Target").
		SetName("MSK").
		SetProvider("AWS").
		SetCategory("Messaging").
		SetTags("kafka","streaming","cloud","managed").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("brokers").
				SetDescription("Set MSK brokers list").
				SetMust(true).
				SetDefault("localhost:9092"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sasl_username").
				SetDescription("Set MSK username").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sasl_password").
				SetDescription("Set MSK password").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("topic").
				SetDescription("Set MSK topic").
				SetMust(true).
				SetDefault(""),
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
