package http

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.http").
		SetDescription("HTTP/Rest Target").
		AddProperty(
			common.NewProperty().
				SetKind("condition").
				SetName("auth_type").
				SetDescription("Sets Auth type").
				SetMust(true).
				SetOptions([]string{"No Auth", "Basic", "Token"}).
				SetDefault("No Auth").
				NewCondition("No Auth", []*common.Property{
					common.NewProperty().
						SetKind("null").
						SetName("auth_type").
						SetDescription("Sets Auth type").
						SetMust(true).
						SetDefault("no_auth"),
				}).
				NewCondition("Basic", []*common.Property{
					common.NewProperty().
						SetKind("null").
						SetName("auth_type").
						SetDescription("Sets Auth type").
						SetMust(true).
						SetDefault("basic"),
					common.NewProperty().
						SetKind("string").
						SetName("username").
						SetDescription("Sets Basic auth username").
						SetMust(true).
						SetDefault(""),
					common.NewProperty().
						SetKind("string").
						SetName("password").
						SetDescription("Sets Basic auth password").
						SetMust(true).
						SetDefault(""),
				}).
				NewCondition("Token", []*common.Property{
					common.NewProperty().
						SetKind("null").
						SetName("auth_type").
						SetDescription("Sets Auth type").
						SetMust(true).
						SetDefault("auth_token"),
					common.NewProperty().
						SetKind("multilines").
						SetName("token").
						SetDescription("Sets Auth token").
						SetMust(true).
						SetDefault(""),
				}),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("proxy").
				SetDescription("Sets Proxy address").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("root_certificate").
				SetDescription("Sets Root Certificate").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("client_private_key").
				SetDescription("Sets Client private key").
				SetMust(false).
				SetDefault(""),
		).AddProperty(
		common.NewProperty().
			SetKind("multilines").
			SetName("client_public_key").
			SetDescription("Sets Client public key").
			SetMust(false).
			SetDefault(""),
	).AddProperty(
		common.NewProperty().
			SetKind("map").
			SetName("default_headers").
			SetDescription("Sets Default headers  (key1=value1;key2=value2;...)").
			SetMust(false).
			SetDefault(""),
	)
}
