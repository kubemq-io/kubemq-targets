package query

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("kubemq.query").
		SetDescription("Kubemq Query Source").
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
				SetDescription("Set how many concurrent Query sources to run").
				SetMust(false).
				SetDefault("1"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("client_id").
				SetDescription("Set Query connection client Id").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("auth_token").
				SetDescription("Set Query connection authentication token").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("auto_reconnect").
				SetDescription("Set Query auto reconnection ").
				SetMust(false).
				SetDefault("true"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("reconnect_interval_seconds").
				SetDescription("Set Query auto reconnection interval in seconds ").
				SetMust(false).
				SetDefault("0"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_reconnects").
				SetDescription("Set Query auto reconnection max reconnects").
				SetMust(false).
				SetDefault("0"),
		)
}
