package redshift

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.redshift.service").
		SetDescription("AWS Redshift Service Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_key").
				SetDescription("Set Redshift Service aws key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_secret_key").
				SetDescription("Set Redshift Service aws secret key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("region").
				SetDescription("Set Redshift Service aws region").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set Redshift Service token").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Redshift Service execution method").
				SetOptions([]string{"create_tags", "delete_tags", "list_tags", "list_snapshots", "list_snapshots_by_tags_keys", "list_snapshots_by_tags_values", "describe_cluster", "list_clusters", "list_clusters_by_tags_keys", "list_clusters_by_tags_values"}).
				SetDefault("create_tags").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("resource_arn").
				SetKind("string").
				SetDescription("Set Redshift Service resource arn").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("resource_name").
				SetKind("string").
				SetDescription("Set Redshift Service resource name").
				SetDefault("").
				SetMust(false),
		)
}
