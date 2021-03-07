package amazonmq

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.amazonmq").
		SetDescription("AWS AmazonMQ Target").
		SetName("AmazonMQ").
		SetProvider("AWS").
		SetCategory("Messaging").
		SetTags("queue","cloud","mq").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("host").
				SetTitle("Host Address").
				SetDescription("Set AmazonMQ host").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Set AmazonMQ username").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Set AmazonMQ password").
				SetMust(true).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("destination").
				SetKind("string").
				SetDescription("Set AmazonMQ destination").
				SetDefault("").
				SetMust(true),
		)
}
