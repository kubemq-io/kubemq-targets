package events

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.cloudwatch.events").
		SetDescription("AWS Cloudwatch Events Target").
		SetName("Cloudwatch Events").
		SetProvider("AWS").
		SetCategory("Observability").
		SetTags("events","cloud").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_key").
				SetDescription("Set Cloudwatch-Events aws key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_secret_key").
				SetDescription("Set Cloudwatch-Events aws secret key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("region").
				SetDescription("Set Cloudwatch-Events aws region").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set Cloudwatch-Events aws token").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Cloudwatch-Events execution method").
				SetOptions([]string{"put_targets", "send_event", "list_buses"}).
				SetDefault("send_event").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("rule").
				SetKind("string").
				SetDescription("Set Cloudwatch-Events rule").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("detail").
				SetKind("string").
				SetDescription("Set Cloudwatch-Events detail").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("detail_type").
				SetKind("string").
				SetDescription("Set Cloudwatch-Events detail type").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("source").
				SetKind("string").
				SetDescription("Set Cloudwatch-Events source").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("limit").
				SetDescription("Set Cloudwatch-Events limit").
				SetMust(false).
				SetDefault("100").
				SetMin(0).
				SetMax(math.MaxInt32),
		)
}
