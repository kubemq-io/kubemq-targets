package queue_stream

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
				SetName("response_channel").
				SetDescription("Set Queue response channel").
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
				SetKind("int").
				SetName("visibility_timeout_seconds").
				SetDescription("Set long to set visibility in seconds to keep the message during processing").
				SetMust(false).
				SetDefault("3600"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("wait_timeout").
				SetDescription("Set how long to wait in seconds for messages during pull of requests").
				SetMust(false).
				SetDefault("3600"),
		)
}
