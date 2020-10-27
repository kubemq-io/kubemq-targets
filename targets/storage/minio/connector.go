package minio

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.storage.minio").
		SetDescription("Minio Storage Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("endpoint").
				SetDescription("Set Minio endpoint address").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("use_ssl").
				SetDescription("Set Minio SSL connection").
				SetMust(false).
				SetDefault("true"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("access_key_id").
				SetDescription("Set Minio access key id").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("secretAccessKey").
				SetDescription("Set Minio secret access key").
				SetMust(false).
				SetDefault(""),
		)
}
