package redis

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("gcp.cache.redis").
		SetDescription("GCP Memory Store Redis Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("url").
				SetDescription("Set Redis url").
				SetMust(true).
				SetDefault("redis://localhost:6379"),
		)
}
