package hazelcast

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("cache.hazelcast").
		SetDescription("hazelcast source properties").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("address").
				SetDescription("Set hazelcast address connection").
				SetMust(true),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Set Username").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Set Password").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("connection_attempt_limit").
				SetDescription("Set connections attempt limit").
				SetMust(false).
				SetDefault("1").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("connection_attempt_period").
				SetDescription("Set connections attempt period in seconds").
				SetMust(false).
				SetDefault("5").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("connection_timeout").
				SetDescription("Set connections attempt timeout in seconds").
				SetMust(false).
				SetDefault("5").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("server_name").
				SetDescription("Set hazelcast server name").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("ssl").
				SetDescription("Set if use ssl").
				SetMust(false).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("condition").
				SetName("ssl").
				SetOptions([]string{"true", "false"}).
				SetDescription("Set ssl conditions").
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
				SetName("map_name").
				SetDescription("Set map name").
				SetMust(true).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("list_name").
				SetDescription("Set list name").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set execution method").
				SetOptions([]string{"get", "set", "delete", "get_list"}).
				SetDefault("get").
				SetMust(true),
		)
}
