package queue_stream

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("kubemq.queue-stream").
		SetDescription("Kubemq Queue Stream Source").
		SetName("KubeMQ Queue Stream").
		SetProvider("").
		SetCategory("Queue").
		AddProperty(
		common.NewProperty().
			SetKind("string").
			SetName("address").
			SetTitle("KubeMQ Address").
			SetDescription("Set Kubemq grpc endpoint address").
			SetMust(true).
			SetDefault("kubemq-cluster-grpc:50000").
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
				SetTitle("Response Channel").
				SetDescription("Set Queue response channel").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sources").
				SetTitle("Concurrent Connections").
				SetDescription("Set how many concurrent Queue sources to run").
				SetMust(false).
				SetDefault("1"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("client_id").
				SetTitle("Client ID").
				SetDescription("Set Queue connection clients Id").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("auth_token").
				SetTitle("Authentication Token").
				SetDescription("Set Queue connection authentication token").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("visibility_timeout_seconds").
				SetTitle("Visibility Timeout (Seconds)").
				SetDescription("Set long to set visibility in seconds to keep the message during processing").
				SetMust(false).
				SetDefault("3600"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("wait_timeout").
				SetTitle("Wait Timeout (Seconds)").
				SetDescription("Set how long to wait in seconds for messages during pull of requests").
				SetMust(false).
				SetDefault("3600"),
		)
}
