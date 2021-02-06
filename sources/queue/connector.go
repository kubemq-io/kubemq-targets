package queue

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("kubemq.queue").
		SetDescription("Kubemq Queue Source").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("address").
				SetDescription("Set Kubemq grpc endpoint address").
				SetMust(true).
				SetDefault("").
				SetLoadedOptions("kubemq-address"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("channel").
				SetDescription("Set Queue channel").
				SetMust(true).
				SetDefaultFromKey("channel.queue"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sources").
				SetDescription("Set how many concurrent Queue sources to run").
				SetMust(false).
				SetDefault("1"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("client_id").
				SetDescription("Set Queue connection clients Id").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("auth_token").
				SetDescription("Set Queue connection authentication token").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("response_channel").
				SetDescription("Set Queue response channel").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("batch_size").
				SetDescription("Set how many messages will pull in one request").
				SetMust(false).
				SetDefault("1"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("wait_timeout").
				SetDescription("Set how long to wait in seconds for messages during pull of requests").
				SetMust(false).
				SetDefault("60"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_requeue").
				SetDescription("Set how many time to requeue a requests do to target error").
				SetMust(true).
				SetDefault("0"),
		)
}
