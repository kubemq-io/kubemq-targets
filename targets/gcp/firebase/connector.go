package firebase

import (
	"github.com/kubemq-hub/builder/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.gcp.firebase").
		SetDescription("GCP Firebase Target").
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
				SetName("auth_client").
				SetDescription("Sets Firebase target is auth client").
				SetMust(false).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("db_client").
				SetDescription("Sets Firebase target is db client").
				SetMust(false).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("db_url").
				SetDescription("Sets Firebase db url").
				SetMust(false).
				SetDefault(""),
		)

}
