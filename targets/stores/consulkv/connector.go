package consulkv

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("stores.consulkv").
		SetDescription("consulkv source properties").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("address").
				SetDescription("Set consulkv address connection").
				SetMust(true),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("wait_time").
				SetDescription("Set wait time milliseconds ").
				SetMust(false).
				SetDefault("36000").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("scheme").
				SetDescription("Set consulkv scheme").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("data_center").
				SetDescription("Set consulkv data center").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set consulkv token").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token_file").
				SetDescription("Set consulkv token_file").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("insecure_skip_verify").
				SetDescription("Set if insecure skip verify").
				SetMust(false).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("tls").
				SetDescription("Set if use tls").
				SetMust(false).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("condition").
				SetName("tls").
				SetOptions([]string{"true", "false"}).
				SetDescription("Set tls conditions").
				SetMust(true).
				SetDefault("false").
				NewCondition("true", []*common.Property{
					common.NewProperty().
						SetKind("multilines").
						SetName("cert_key").
						SetDescription("Set certificate key").
						SetMust(false).
						SetDefault(""),
					common.NewProperty().
						SetKind("multilines").
						SetName("cert_file").
						SetDescription("Set certificate file").
						SetMust(false).
						SetDefault(""),
				}),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("key").
				SetDescription("Set key").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("data_center").
				SetDescription("Set data center").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("near").
				SetDescription("Set near").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("filter").
				SetDescription("Set filter").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("prefix").
				SetDescription("Set prefix").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("bool").
				SetName("allow_stale").
				SetDescription("Set if allow stale").
				SetMust(false).
				SetDefault("false"),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("bool").
				SetName("require_consistent").
				SetDescription("Set if require consistent").
				SetMust(false).
				SetDefault("false"),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("bool").
				SetName("user_cache").
				SetDescription("Set if use user cache").
				SetMust(false).
				SetDefault("false"),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("max_age").
				SetDescription("Set max age milliseconds ").
				SetMust(false).
				SetDefault("36000").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("stale_if_error").
				SetDescription("Set stale if error time in milliseconds").
				SetMust(false).
				SetDefault("36000").
				SetMin(1).
				SetMax(math.MaxInt32),
		)
}
