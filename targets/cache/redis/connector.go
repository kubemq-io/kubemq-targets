package redis

import (
	"github.com/kubemq-hub/builder/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.cache.redis").
		SetDescription("Redis Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("url").
				SetDescription("Sets Redis url").
				SetMust(true).
				SetDefault("redis://localhost:6379"),
		)
}
