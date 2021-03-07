package elasticsearch

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.elasticsearch").
		SetDescription("AWS Elastic Search Target").
		SetName("Elasticsearch").
		SetProvider("AWS").
		SetCategory("Store").
		SetTags("db","log","cloud","managed").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_key").
				SetDescription("Set Elastic Search aws key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_secret_key").
				SetDescription("Set Elastic Search aws secret key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set Elastic Search aws token").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Elastic Search execution method").
				SetOptions([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}).
				SetDefault("GET").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("region").
				SetKind("string").
				SetDescription("Set Elastic Search region").
				SetDefault("").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("domain").
				SetKind("string").
				SetDescription("Set Elastic Search domain").
				SetDefault("").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("index").
				SetDescription("Set Elastic Search index").
				SetMust(true).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("endpoint").
				SetDescription("Set Elastic Search endpoint").
				SetMust(true).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("id").
				SetDescription("Set Elastic Search id").
				SetMust(true).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("multilines").
				SetName("json").
				SetDescription("Set Elastic Search json").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("service").
				SetDescription("Set Elastic Search service").
				SetMust(false).
				SetDefault("es"),
		)
}
