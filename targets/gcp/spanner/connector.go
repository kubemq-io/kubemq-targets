package spanner

import (
	"github.com/kubemq-hub/builder/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.gcp.spanner").
		SetDescription("GCP Spanner Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("db").
				SetDescription("Sets GCP Spanner DB").
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
