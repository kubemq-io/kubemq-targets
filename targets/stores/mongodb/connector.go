package mongodb

import (
	"github.com/kubemq-hub/builder/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.stores.mongodb").
		SetDescription("MongoDB Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("host").
				SetDescription("Sets MongoDB host address").
				SetMust(true).
				SetDefault("localhost:27017"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Sets MongoDB username").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Sets MongoDB password").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("collection").
				SetDescription("Sets MongoDB collection").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("params").
				SetDescription("Sets MongoDB params").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("read_concurrency").
				SetDescription("Sets MongoDB read concurrency").
				SetOptions([]string{"local", "majority", "available", "linearizable", "snapshot"}).
				SetMust(false).
				SetDefault("local"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("write_concurrency").
				SetDescription("Sets MongoDB write concurrency").
				SetOptions([]string{"majority", "Other"}).
				SetMust(false).
				SetDefault("majority"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("operation_timeout_seconds").
				SetDescription("Sets MongoDB operation timeout seconds").
				SetMust(false).
				SetDefault("30").
				SetMin(0).
				SetMax(math.MaxInt32),
		)
}
