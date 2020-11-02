package pubsub

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("gcp.pubsub").
		SetDescription("GCP PubSub Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("project_id").
				SetDescription("Set GCP project ID").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("credentials").
				SetDescription("Set GCP credentials").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("retries").
				SetDescription("Set PubSub sending message retries").
				SetMust(false).
				SetDefault("0").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("topic_id").
				SetKind("string").
				SetDescription("Set PubSub request topic id").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("tags").
				SetKind("string").
				SetDescription("Set PubSub request tags").
				SetDefault("").
				SetMust(false),
		)
}
