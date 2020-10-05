package openfaas

import (
	"github.com/kubemq-hub/builder/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.serverless.openfaas").
		SetDescription("Openfaas Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("gateway").
				SetDescription("Sets Openfaas gateway address").
				SetMust(true).
				SetDefault("localhost:27017"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Sets Openfaas username").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Sets Openfaas password").
				SetMust(true).
				SetDefault(""),
		)
}
