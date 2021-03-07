package storage

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("gcp.storage").
		SetDescription("GCP Storage Target").
		SetName("Storage").
		SetProvider("GCP").
		SetCategory("Storage").
		SetTags("db","filesystem","object","cloud","managed").
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("credentials").
				SetTitle("Json Credentials").
				SetDescription("Set GCP credentials").
				SetMust(true).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set GCP Storage method").
				SetOptions([]string{"upload", "create_bucket", "download", "delete", "rename", "copy", "move", "list"}).
				SetDefault("create_bucket").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("object").
				SetKind("string").
				SetDescription("Set object name to save the file under").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("rename_object").
				SetKind("string").
				SetDescription("Set GCP name to change the file name").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("bucket").
				SetKind("string").
				SetDescription("Set Storage bucket name").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("dst_bucket").
				SetKind("string").
				SetDescription("Set the bucket name of the destination").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("path").
				SetKind("string").
				SetDescription("Set path to the file for upload").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("project_id").
				SetKind("string").
				SetDescription("Set GCP storage project id ").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("storage_class").
				SetKind("string").
				SetDescription("Set GCP storage class").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("location").
				SetKind("string").
				SetDescription("Set GCP storage location").
				SetDefault("").
				SetMust(false),
		)
}
