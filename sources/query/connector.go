package query

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("kubemq.query").
		SetDescription("Kubemq Query Source").
		SetName("KubeMQ Query").
		SetProvider("").
		SetCategory("RPC").
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
				SetDescription("Set Query channel").
				SetMust(true).
				SetDefaultFromKey("channel.query"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("group").
				SetDescription("Set Query channel group").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sources").
				SetTitle("Concurrent Connections").
				SetDescription("Set how many concurrent Query sources to run").
				SetMust(false).
				SetDefault("1"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("client_id").
				SetTitle("Client ID").
				SetDescription("Set Query connection client Id").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("auth_token").
				SetTitle("Authentication Token").
				SetDescription("Set Query connection authentication token").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("auto_reconnect").
				SetTitle("Reconnect Automatically").
				SetDescription("Set auto reconnection ").
				SetMust(false).
				SetDefault("true"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("reconnect_interval_seconds").
				SetTitle("Reconnection Interval (Seconds)").
				SetDescription("Set auto reconnection interval in seconds ").
				SetMust(false).
				SetDefault("0"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_reconnects").
				SetTitle("Max Reconnections").
				SetDescription("Set auto reconnection max reconnects").
				SetMust(false).
				SetDefault("0"),
		)
}
