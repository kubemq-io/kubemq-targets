package bigquery

import "github.com/kubemq-hub/builder/common"

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.gcp.bigquery").
		SetDescription("GCP Bigquery Target").
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
		)
}
