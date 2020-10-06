package memcached

import (
	"github.com/kubemq-hub/builder/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.gcp.cache.memcached").
		SetDescription("GCP Memory Store Memcached Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("hosts").
				SetDescription("Sets Memcached hosts").
				SetMust(true).
				SetDefault("localhost:11211"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_idle_connections").
				SetDescription("Sets Memcached max idle connections").
				SetDefault("2").
				SetMin(1).
				SetMax(math.MaxInt32).
				SetMust(false),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("default_timeout_seconds").
				SetDescription("Sets Memcached default timeout seconds").
				SetDefault("30").
				SetMin(1).
				SetMax(math.MaxInt32).
				SetMust(false),
		)
}
