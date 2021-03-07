package s3

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.s3").
		SetDescription("AWS S3 Target").
		SetName("S3").
		SetProvider("AWS").
		SetCategory("Storage").
		SetTags("filesystem","object","db","cloud","managed").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_key").
				SetDescription("Set S3 aws key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_secret_key").
				SetDescription("Set S3 aws secret key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("region").
				SetDescription("Set S3 aws region").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set S3 token").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("downloader").
				SetDescription("Create S3 downloader instance").
				SetMust(false).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("uploader").
				SetDescription("Create S3 uploader instance").
				SetMust(false).
				SetDefault("false"),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set S3 execution method").
				SetOptions([]string{"list_buckets", "list_bucket_items", "create_bucket", "delete_bucket", "delete_item_from_bucket", "delete_all_items_from_bucket", "upload_item", "copy_item", "get_item"}).
				SetDefault("upload_item").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("bucket_name").
				SetKind("string").
				SetDescription("Set S3 bucket name").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("copy_source").
				SetKind("string").
				SetDescription("Set S3 copy source").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("bool").
				SetName("wait_for_completion").
				SetDescription("Set S3 wait for completion until retuning response").
				SetMust(false).
				SetDefault("false"),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("item_name").
				SetDescription("Set S3 item name").
				SetMust(false).
				SetDefault(""),
		)

}
