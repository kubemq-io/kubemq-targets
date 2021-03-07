package mqtt

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("messaging.mqtt").
		SetDescription("MQTT Messaging Target").
		SetName("MQTT").
		SetProvider("").
		SetCategory("Messaging").
		SetTags("iot","pub/sub").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("host").
				SetTitle("Host Address").
				SetDescription("Set MQTT broker host").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Set MQTT broker username").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Set MQTT broker password").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("client_id").
				SetTitle("Client ID").
				SetDescription("Set MQTT broker client id").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("topic").
				SetDescription("Set MQTT topic").
				SetDefault("").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("qos").
				SetDescription("Set MQTT qos level").
				SetMust(true).
				SetMin(0).
				SetMax(2).
				SetDefault("0"),
		)
}
