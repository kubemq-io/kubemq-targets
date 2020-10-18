package redis

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.gcp.cache.redis").
		SetDescription("GCP Memory Store Redis Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("url").
				SetDescription("Sets Redis url").
				SetMust(true).
				SetDefault("redis://localhost:6379"),
		)
}
