package logs

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.cloudwatch.logs").
		SetDescription("AWS Cloudwatch Logs Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_key").
				SetDescription("Set Cloudwatch-Logs aws key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_secret_key").
				SetDescription("Set Cloudwatch-Logs aws secret key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("region").
				SetDescription("Set Cloudwatch-Logs aws region").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set Cloudwatch-Logs aws token").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Cloudwatch-Logs execution method").
				SetOptions([]string{"create_log_event_stream", "describe_log_event_stream", "delete_log_event_stream", "put_log_event", "get_log_event", "create_log_group", "delete_log_group", "describe_log_group"}).
				SetDefault("put_log_event").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("sequence_token").
				SetKind("string").
				SetDescription("Set Cloudwatch-Logs sequence token").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("log_group_name").
				SetKind("string").
				SetDescription("Set Cloudwatch-Logs log group name").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("log_stream_name").
				SetKind("string").
				SetDescription("Set Cloudwatch-Logs sequence log stream name").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("log_group_prefix").
				SetKind("string").
				SetDescription("Set Cloudwatch-Logs log group prefix").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("limit").
				SetDescription("Set Cloudwatch-Logs limit").
				SetMust(false).
				SetDefault("100").
				SetMin(0).
				SetMax(math.MaxInt32),
		)

}
