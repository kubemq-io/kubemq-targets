package filesystem

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("storage.filesystem").
		SetDescription("Local Filesystem Target").
		SetName("File System").
		SetProvider("").
		SetCategory("Storage").
		SetTags("filesystem").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("base_path").
				SetTitle("Destination Path").
				SetDescription("Set local file system base path").
				SetMust(true).
				SetDefault("./"),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set file system method").
				SetOptions([]string{"save", "load", "delete", "list"}).
				SetDefault("").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("path").
				SetKind("string").
				SetDescription("Set path").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("filename").
				SetKind("string").
				SetDescription("Set filename").
				SetDefault("").
				SetMust(true),
		)
}
