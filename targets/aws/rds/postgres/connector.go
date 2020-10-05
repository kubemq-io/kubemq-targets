package command

import (
	"github.com/kubemq-hub/builder/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("source.command").
		SetDescription("Kubemq Command Source").
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
				SetDescription("Sets Command channel").
				SetMust(true).
				SetDefault("commands"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("group").
				SetDescription("Sets Command channel group").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("client_id").
				SetDescription("Sets Command connection client Id").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("auth_token").
				SetDescription("Sets Command connection authentication token").
				SetMust(false).
				SetDefault(""),
		)
}
