package bigtable

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.gcp.bigtable").
		SetDescription("GCP Bigtable Target").
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
				SetKind("string").
				SetName("instance").
				SetDescription("Sets Bigtable instance").
				SetMust(true).
				SetDefault(""),
		)
}
