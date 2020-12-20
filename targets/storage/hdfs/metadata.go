package s3

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
	"os"
)

type metadata struct {
	method string

	filePath    string
	oldFilePath string
	fileMode    os.FileMode
}

var methodsMap = map[string]string{
	"write_file":  "write_file",
	"remove_file": "remove_file",
	"read_file":   "read_file",
	"rename_file": "rename_file",
	"mkdir":       "mkdir",
	"stat":        "stat",
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}
	m.filePath, err = meta.MustParseString("file_path")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing file_path, %w", err)
	}
	if m.method == "rename_file" {
		m.oldFilePath, err = meta.MustParseString("old_file_path")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing old_file_path, %w", err)
		}
	}
	m.fileMode = meta.ParseOSFileMode("file_mode", os.FileMode(0755))
	return m, nil
}
