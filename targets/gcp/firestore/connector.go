package firestore

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("gcp.firestore").
		SetDescription("GCP Firestore Target").
		SetName("Firestore").
		SetProvider("GCP").
		SetCategory("Store").
		SetTags("db","no-sql","cloud","managed").
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
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set GCP Firestore execution method").
				SetOptions([]string{"documents_all", "document_key", "delete_document_key", "add"}).
				SetDefault("add").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("collection").
				SetKind("string").
				SetDescription("Set Firestore collection name").
				SetDefault("").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("item").
				SetKind("string").
				SetDescription("Set Firestore item name").
				SetDefault("").
				SetMust(false),
		)
}
