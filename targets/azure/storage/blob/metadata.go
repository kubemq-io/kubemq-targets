package blob

import (
	"fmt"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/kubemq-io/kubemq-targets/types"
)

const (
	DeleteSnapshotsOptionInclude = "include"
	DeleteSnapshotsOptionNone    = ""
	DeleteSnapshotsOptionOnly    = "only"

	DefaultRetryRequests = 1
	DefaultBlockSize     = 4194304
	DefaultParallelism   = 16

	DefaultCount  = 0
	DefaultOffset = 0
)

var methodsMap = map[string]string{
	"upload": "upload",
	"get":    "get",
	"delete": "delete",
}

var deleteSnapShotTypes = map[string]string{
	"include": "include",
	"only":    "only",
	"":        "",
}

type metadata struct {
	method                    string
	fileName                  string
	serviceUrl                string
	blockSize                 int64
	parallelism               uint16
	offset                    int64
	count                     int64
	deleteSnapshotsOptionType azblob.DeleteSnapshotsOptionType
	maxRetryRequests          int
	blobMetadata              map[string]string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}
	if m.method == "delete" {
		deleteSnapshotsOptionType, err := meta.ParseStringMap("delete_snapshots_option_type", deleteSnapShotTypes)
		if err != nil {
			return metadata{}, meta.GetValidSupportedTypes(deleteSnapShotTypes, "delete_snapshots_option_type")
		}
		switch deleteSnapshotsOptionType {
		case DeleteSnapshotsOptionInclude:
			m.deleteSnapshotsOptionType = azblob.DeleteSnapshotsOptionInclude
		case DeleteSnapshotsOptionOnly:
			m.deleteSnapshotsOptionType = azblob.DeleteSnapshotsOptionOnly
		case DeleteSnapshotsOptionNone:
			m.deleteSnapshotsOptionType = azblob.DeleteSnapshotsOptionNone
		}
	}
	m.fileName, err = meta.MustParseString("file_name")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing file_name , %w", err)
	}
	m.serviceUrl, err = meta.MustParseString("service_url")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing service_url , %w", err)
	}

	blobMetadata, err := meta.MustParseJsonMap("blob_metadata")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing blob_metadata, %w", err)
	} else {
		m.blobMetadata = blobMetadata
	}

	m.blockSize = int64(meta.ParseInt("block_size", DefaultBlockSize))
	m.parallelism = uint16(meta.ParseInt("parallelism", DefaultParallelism))
	m.count = int64(meta.ParseInt("count", DefaultCount))
	m.offset = int64(meta.ParseInt("offset", DefaultOffset))
	m.maxRetryRequests = meta.ParseInt("max_retry_request", DefaultRetryRequests)
	return m, nil
}
