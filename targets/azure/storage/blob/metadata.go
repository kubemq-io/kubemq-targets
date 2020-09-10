package blob

import (
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/kubemq-hub/kubemq-targets/types"
)

const (
	PublicAccessBlob      = "blob"
	PublicAccessContainer = "container"
	DefaultBlockSize      = 4194304
	DefaultParallelism    = 16
)

var publicAccessType = map[string]string{
	"blob":      "blob",
	"container": "container",
}

var methodsMap = map[string]string{
	"upload": "upload",
	"get":    "get",
}

type metadata struct {
	method                    string
	accessType                azblob.PublicAccessType
	fileName                  string
	serviceUrl                string
	blockSize                 int64
	parallelism               uint16
	offset                    int64
	count                     int64
	deleteSnapshotsOptionType string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}
	accessType, err := meta.ParseStringMap("publicAccessType", publicAccessType)
	if err != nil {
		m.accessType = azblob.PublicAccessNone
	}
	if accessType == PublicAccessBlob {
		m.accessType = azblob.PublicAccessBlob
	} else if accessType == PublicAccessContainer {
		m.accessType = azblob.PublicAccessContainer
	}
	m.blockSize = int64(meta.ParseInt("block_size", DefaultBlockSize))
	m.blockSize = int64(meta.ParseInt("parallelism", DefaultParallelism))
	if m.method == "upload" || m.method == "get" {
		m.fileName, err = meta.MustParseString("file_name")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing file_name , %w", err)
		}
		m.serviceUrl, err = meta.MustParseString("service_url")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing service_url , %w", err)
		}
	}
	m.deleteSnapshotsOptionType = meta.ParseString("delete_snapshots_option_type")

	return m, nil
}
