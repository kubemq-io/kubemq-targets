package http

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("http").
		SetDescription("HTTP/Rest Target").
		AddProperty(
			common.NewProperty().
				SetKind("condition").
				SetName("auth_type").
				SetDescription("Set Auth type").
				SetMust(true).
				SetOptions([]string{"No Auth", "Basic", "Token"}).
				SetDefault("No Auth").
				NewCondition("No Auth", []*common.Property{
					common.NewProperty().
						SetKind("null").
						SetName("auth_type").
						SetDescription("Set Auth type").
						SetMust(true).
						SetDefault("no_auth"),
				}).
				NewCondition("Basic", []*common.Property{
					common.NewProperty().
						SetKind("null").
						SetName("auth_type").
						SetDescription("Set Auth type").
						SetMust(true).
						SetDefault("basic"),
					common.NewProperty().
						SetKind("string").
						SetName("username").
						SetDescription("Set Basic auth username").
						SetMust(true).
						SetDefault(""),
					common.NewProperty().
						SetKind("string").
						SetName("password").
						SetDescription("Set Basic auth password").
						SetMust(true).
						SetDefault(""),
				}).
				NewCondition("Token", []*common.Property{
					common.NewProperty().
						SetKind("null").
						SetName("auth_type").
						SetDescription("Set Auth type").
						SetMust(true).
						SetDefault("auth_token"),
					common.NewProperty().
						SetKind("multilines").
						SetName("token").
						SetDescription("Set Auth token").
						SetMust(true).
						SetDefault(""),
				}),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("proxy").
				SetDescription("Set Proxy address").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("root_certificate").
				SetDescription("Set Root Certificate").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("client_private_key").
				SetDescription("Set Client private key").
				SetMust(false).
				SetDefault(""),
		).AddProperty(
		common.NewProperty().
			SetKind("multilines").
			SetName("client_public_key").
			SetDescription("Set Client public key").
			SetMust(false).
			SetDefault(""),
	).AddProperty(
		common.NewProperty().
			SetKind("map").
			SetName("default_headers").
			SetDescription("Set Default headers  (key1=value1;key2=value2;...)").
			SetMust(false).
			SetDefault(""),
	)
}
