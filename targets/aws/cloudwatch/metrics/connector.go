package metrics

import "github.com/kubemq-hub/builder/connector/common"

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.cloudwatch.metrics").
		SetDescription("AWS Cloudwatch Metrics Target").
		SetName("Cloudwatch Metrics").
		SetProvider("AWS").
		SetCategory("Observability").
		SetTags("metrics", "cloud").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_key").
				SetDescription("Set Cloudwatch-Metrics aws key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_secret_key").
				SetDescription("Set Cloudwatch-Metrics aws secret key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("region").
				SetDescription("Set Cloudwatch-Metrics aws region").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set Cloudwatch-Metrics aws token").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Cloudwatch-Metrics execution method").
				SetOptions([]string{"put_metrics", "list_metrics"}).
				SetDefault("put_metrics").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("namespace").
				SetKind("string").
				SetDescription("Set Cloudwatch-Metrics namespace").
				SetDefault("").
				SetMust(false),
		)
}
