package sqs

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.aws.sqs").
		SetDescription("AWS SQS Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_key").
				SetDescription("Sets SQS aws key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_secret_key").
				SetDescription("Sets SQS aws secret key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("region").
				SetDescription("Sets SQS aws region").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Sets SQS token").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_receive").
				SetDescription("Sets SQS max receive").
				SetMust(false).
				SetDefault("0").
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("retries").
				SetDescription("Sets SQS number of retries on failed send request").
				SetMust(false).
				SetDefault("0").
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("default_delay").
				SetDescription("Sets SQS default delay in seconds").
				SetMust(false).
				SetDefault("10").
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("dead_letter").
				SetDescription("Sets SQS dead letter queue").
				SetMust(false).
				SetDefault(""),
		)
}
