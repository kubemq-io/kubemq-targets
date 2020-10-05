package couchbase

import (
	"github.com/kubemq-hub/builder/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.stores.couchbase").
		SetDescription("Couchbase Target").
		//
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("host").
				SetDescription("Sets Couchbase url").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Sets Couchbase username").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Sets Couchbase password").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("bucket").
				SetDescription("Sets Couchbase bucket").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("collection").
				SetDescription("Sets Couchbase collection").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("num_to_replicate").
				SetDescription("Sets Couchbase number of nodes to replicate").
				SetMust(false).
				SetDefault("1").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("num_to_persist").
				SetDescription("Sets Couchbase number of node to persist").
				SetMust(false).
				SetDefault("1").
				SetMin(1).
				SetMax(math.MaxInt32),
		)
}
