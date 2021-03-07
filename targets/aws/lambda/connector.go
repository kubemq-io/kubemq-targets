package lambda

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.lambda").
		SetDescription("AWS Lambda Target").
		SetName("Lambda").
		SetProvider("AWS").
		SetCategory("Serverless").
		SetTags("faas","cloud","managed").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_key").
				SetDescription("Set Lambda aws key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_secret_key").
				SetDescription("Set Lambda aws secret key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("region").
				SetDescription("Set Lambda aws region").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set Lambda token").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Lambda execution method").
				SetOptions([]string{"list", "create", "run", "delete"}).
				SetDefault("run").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("zip_file_name").
				SetKind("string").
				SetDescription("Set Lambda zip file name").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("handler_name").
				SetKind("string").
				SetDescription("Set Lambda handler name").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("role").
				SetKind("string").
				SetDescription("Set Lambda role").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("runtime").
				SetKind("string").
				SetDescription("Set Lambda runtime version").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("function_name").
				SetKind("string").
				SetDescription("Set Lambda function name").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("description").
				SetKind("string").
				SetDescription("Set Lambda description").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("memory_size").
				SetDescription("Set Lambda memory size").
				SetMust(false).
				SetDefault("256").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("int").
				SetName("timeout").
				SetDescription("Set Lambda timeout").
				SetMust(false).
				SetDefault("15").
				SetMin(0).
				SetMax(math.MaxInt32),
		)
}
