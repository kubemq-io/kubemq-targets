package msk

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.aws.msk").
		SetDescription("AWS MSK Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("brokers").
				SetDescription("Sets MSK brokers list").
				SetMust(true).
				SetDefault("localhost:9092"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sasl_username").
				SetDescription("Sets MSK username").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sasl_password").
				SetDescription("Sets MSK password").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("topic").
				SetDescription("Sets MSK topic").
				SetMust(true).
				SetDefault(""),
		)

}
