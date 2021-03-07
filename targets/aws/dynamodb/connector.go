package dynamodb

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.dynamodb").
		SetDescription("AWS Dynamodb Target").
		SetName("DynamoDB").
		SetProvider("AWS").
		SetCategory("Store").
		SetTags("db","no-sql","cloud","managed").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_key").
				SetDescription("Set Dynamodb aws key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_secret_key").
				SetDescription("Set Dynamodb aws secret key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("region").
				SetDescription("Set Dynamodb aws region").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set Dynamodb aws token").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Dynamodb execution method").
				SetOptions([]string{"list_tables", "create_table", "delete_table", "insert_item", "get_item", "delete_item", "update_item"}).
				SetDefault("insert_item").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("table_name").
				SetKind("string").
				SetDescription("Set Dynamodb table name").
				SetDefault("").
				SetMust(false),
		)
}
