package minio

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("storage.minio").
		SetDescription("Minio Storage Target").
		SetName("Minio").
		SetProvider("").
		SetCategory("Storage").
		SetTags("s3").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("endpoint").
				SetTitle("Endpoint").
				SetDescription("Set Minio endpoint address").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("use_ssl").
				SetTitle("USE SSL").
				SetDescription("Set Minio SSL connection").
				SetMust(false).
				SetDefault("true"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("access_key_id").
				SetTitle("Access Key ID").
				SetDescription("Set Minio access key id").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("secret_access_key").
				SetTitle("Access Key Secret").
				SetDescription("Set Minio secret access key").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Minio method").
				SetOptions([]string{"make_bucket", "list_buckets", "bucket_exists", "remove_bucket", "list_objects", "put", "get", "remove"}).
				SetDefault("make_bucket").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("param1").
				SetKind("string").
				SetDescription("Set Minio bucket name").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("param2").
				SetKind("string").
				SetDescription("Set Minio object name").
				SetDefault("").
				SetMust(false),
		)
}
