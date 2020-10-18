package mqtt

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.messaging.mqtt").
		SetDescription("MQTT Messaging Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("host").
				SetDescription("Sets MQTT broker host").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Sets MQTT broker username").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Sets MQTT broker password").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("client_id").
				SetDescription("Sets MQTT broker client id").
				SetMust(true).
				SetDefault(""),
		)
}
