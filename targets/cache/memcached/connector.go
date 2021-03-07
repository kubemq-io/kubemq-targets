package memcached

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("cache.memcached").
		SetDescription("Memcached Target").
		SetName("Memcached").
		SetProvider("").
		SetCategory("Cache").
		SetTags("db").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("hosts").
				SetTitle("Hosts Address").
				SetDescription("Set Memcached hosts").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_idle_connections").
				SetDescription("Set Memcached max idle connections").
				SetDefault("2").
				SetMin(1).
				SetMax(math.MaxInt32).
				SetMust(false),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("default_timeout_seconds").
				SetTitle("Default Timeout (Seconds)").
				SetDescription("Set Memcached default timeout seconds").
				SetDefault("30").
				SetMin(1).
				SetMax(math.MaxInt32).
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Memcached execution method").
				SetOptions([]string{"get", "set", "delete"}).
				SetDefault("get").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("key").
				SetKind("string").
				SetDescription("Set Memcached key").
				SetMust(true),
		)
}
