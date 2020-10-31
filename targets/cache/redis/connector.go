package redis

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
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
				SetDefault("redis://redis.redis:6379"),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Select Redis method").
				SetOptions([]string{"get","set","delete"}).
				SetDefault("get").
				SetMust(true),
					).
		AddMetadata(
			common.NewMetadata().
				SetName("key").
				SetKind("string").
				SetDescription("Set Redis key").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("etag").
				SetKind("int").
				SetDescription("Set Redis etag").
				SetDefault("0").
				SetMin(0).
				SetMax(math.MaxInt16).
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("concurrency").
				SetKind("string").
				SetDescription("Set Redis write concurrency").
				SetOptions([]string{"first-write","last-write",""}).
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("consistency").
				SetKind("string").
				SetDescription("Set Redis read consistency").
				SetOptions([]string{"strong","eventual",""}).
				SetDefault("").
				SetMust(false),
	)


}
