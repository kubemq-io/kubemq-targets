package firebase

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("gcp.firebase").
		SetDescription("GCP Firebase Target").
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
				SetKind("bool").
				SetName("auth_client").
				SetDescription("Set Firebase target is a auth client").
				SetMust(false).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("db_client").
				SetDescription("Set Firebase target is a db client").
				SetMust(false).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("db_url").
				SetDescription("Set Firebase db url").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("messaging_client").
				SetDescription("Set Firebase target is a messaging client").
				SetMust(false).
				SetDefault("false"),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set GCP Firebase execution method").
				SetOptions([]string{"custom_token", "verify_token", "retrieve_user", "create_user", "delete_user", "delete_multiple_users", "list_users", "get_db", "delete_db", "set_db", "send_message", "send_multi"}).
				SetDefault("create_user").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("retrieve_by").
				SetKind("string").
				SetDescription("Set GCP Firebase retrieve by type").
				SetOptions([]string{"by_uid", "by_email", "by_phone"}).
				SetDefault("by_uid").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("token_id").
				SetKind("string").
				SetDescription("Set Firebase token id").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("uid").
				SetKind("string").
				SetDescription("Set Firebase user uid").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("email").
				SetKind("string").
				SetDescription("Set Firebase user email").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("phone").
				SetKind("string").
				SetDescription("Set Firebase user phone").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("ref_path").
				SetKind("string").
				SetDescription("Set Firebase reference path").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("child_ref").
				SetKind("string").
				SetDescription("Set Firebase child path").
				SetDefault("").
				SetMust(false),
		)
}
