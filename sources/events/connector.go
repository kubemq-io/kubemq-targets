package events

import (
	"github.com/kubemq-hub/builder/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("source.events").
		SetDescription("Kubemq Events Source").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("address").
				SetDescription("Sets Kubemq grpc endpoint address").
				SetMust(true).
				SetDefault("localhost:50000"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("channel").
				SetDescription("Sets Events channel").
				SetMust(true).
				SetDefault("events"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("group").
				SetDescription("Sets Events channel group").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("response_channel").
				SetDescription("Sets Events response channel").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("client_id").
				SetDescription("Sets Events connection client Id").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("auth_token").
				SetDescription("Sets Events connection authentication token").
				SetMust(false).
				SetDefault(""),
		)
}
