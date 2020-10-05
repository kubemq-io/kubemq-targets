package mysql

import (
	"github.com/kubemq-hub/builder/common"
)

// TODO
func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.stores.mysql").
		SetDescription("MySQL Target")
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
