package redshift

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.aws.rds.redshift").
		SetDescription("AWS RDS Redshift Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("connection").
				SetDescription("Sets Redshift connection string").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_idle_connections").
				SetDescription("Sets Redshift max idle connections").
				SetMust(false).
				SetDefault("10").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_open_connections").
				SetDescription("Sets Redshift max open connections").
				SetMust(false).
				SetDefault("100").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("connection_max_lifetime_seconds").
				SetDescription("Sets Redshift connection max lifetime seconds").
				SetMust(false).
				SetDefault("3600").
				SetMin(1).
				SetMax(math.MaxInt32),
		)
	//
	//AddProperty(
	//	common.NewProperty().
	//		SetKind("string").
	//		SetName("address").
	//		SetDescription("Sets Kubemq grpc endpoint address").
	//		SetMust(true).
	//		SetDefault("localhost:50000"),
	//)
}
