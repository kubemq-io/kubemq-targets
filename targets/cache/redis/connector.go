package redis

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("cache.redis").
		SetDescription("Redis Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("url").
				SetDescription("Set Redis url").
				SetMust(true).
				SetDefault("redis://localhost:6379"),
		)
}
