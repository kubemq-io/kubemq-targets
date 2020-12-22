package hdfs

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("storage.hdfs").
		SetDescription("Hadoop Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("address").
				SetDescription("Set Hadoop address").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("user").
				SetDescription("Set Hadoop user").
				SetMust(false).
				SetDefault(""),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Hadoop execution method").
				SetOptions([]string{"write_file", "remove_file", "read_file", "rename_file", "mkdir", "stat"}).
				SetDefault("read_file").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("file_path").
				SetKind("string").
				SetDescription("Set Hadoop file path").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("old_file_path").
				SetKind("string").
				SetDescription("Set Hadoop old file path").
				SetDefault("").
				SetMust(false),
		).
		AddMetadata(
			common.NewMetadata().
				SetKind("string").
				SetName("file_mode").
				SetDescription("Set os file mode").
				SetMust(false).
				SetDefault("0777"),
		)

}
