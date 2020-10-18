package cloudfunctions

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.gcp.cloudfunctions").
		SetDescription("GCP Cloud Functions Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("project_id").
				SetDescription("Sets GCP project ID").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("credentials").
				SetDescription("Sets GCP credentials").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("location_match").
				SetDescription("Sets Cloud Functions location match").
				SetMust(false).
				SetDefault("true"),
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
