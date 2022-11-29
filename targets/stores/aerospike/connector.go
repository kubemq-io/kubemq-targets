package aerospike

import (
	"math"

	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("stores.aerospike").
		SetDescription("Aerospike Target").
		SetName("Aerospike").
		SetProvider("").
		SetCategory("Store").
		SetTags("db").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("host").
				SetTitle("Host Address").
				SetDescription("Set Aerospike host address").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("port").
				SetDescription("Set Aerospike port address").
				SetMust(true).
				SetMin(0).
				SetMax(65355).
				SetDefault("3000"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Set Aerospike username").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Set Aerospike password").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("timeout").
				SetDescription("Set aerospike timeout seconds").
				SetMust(false).
				SetDefault("5").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set aerospike execution method").
				SetOptions([]string{"get", "set", "delete", "get_batch"}).
				SetDefault("get").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("key").
				SetKind("string").
				SetDescription("Set aerospike key").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("user_key").
				SetKind("string").
				SetDescription("Set aerospike user key").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("namespace").
				SetKind("string").
				SetDescription("Set aerospike namespace").
				SetDefault("").
				SetMust(false),
		)
}
