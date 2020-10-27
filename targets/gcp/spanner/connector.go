package spanner

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.gcp.spanner").
		SetDescription("GCP Spanner Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("db").
				SetDescription("Set GCP Spanner DB").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("credentials").
				SetDescription("Set GCP credentials").
				SetMust(true).
				SetDefault(""),
		)
}
