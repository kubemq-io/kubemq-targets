package sns

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.sns").
		SetDescription("AWS SNS Target").
		SetName("SNS").
		SetProvider("AWS").
		SetCategory("Messaging").
		SetTags("pub/sub","cloud","managed").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_key").
				SetDescription("Set SNS aws key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_secret_key").
				SetDescription("Set SNS aws secret key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("region").
				SetDescription("Set SNS aws region").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set SNS token").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set SNS execution method").
				SetOptions([]string{"list_topics", "list_subscriptions", "list_subscriptions_by_topic", "create_topic", "subscribe", "send_message", "delete_topic"}).
				SetDefault("send_message").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("topic").
				SetKind("string").
				SetDescription("Set SNS topic").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("end_point").
				SetKind("string").
				SetDescription("Set SNS end point").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("bool").
				SetName("return_subscription").
				SetDescription("Set SNS return subscription").
				SetMust(false).
				SetDefault("false"),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("target_arn").
				SetDescription("Set SNS target arn").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("message").
				SetDescription("Set SNS message").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("phone_number").
				SetDescription("Set SNS phone number").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("subject").
				SetDescription("Set SNS subject").
				SetMust(false).
				SetDefault(""),
		)
}
