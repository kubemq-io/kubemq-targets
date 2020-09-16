package files

import (
	"fmt"
	"github.com/Azure/azure-storage-file-go/azfile"
	"github.com/kubemq-hub/kubemq-targets/types"
)

const (
	DeleteSnapshotsOptionInclude = "include"
	DeleteSnapshotsOptionNone    = ""

	DefaultRetryRequests = 0
	DefaultRangeSize     = 4194304
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
	"":        "",
}

type metadata struct {
	method                    string
	fileName                  string
	serviceUrl                string
	rangeSize                 int64
	parallelism               uint16
	offset                    int64
	count                     int64
	deleteSnapshotsOptionType azfile.DeleteSnapshotsOptionType
	maxRetryRequests          int
	fileMetadata              map[string]string
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
			m.deleteSnapshotsOptionType = azfile.DeleteSnapshotsOptionInclude
		case DeleteSnapshotsOptionNone:
			m.deleteSnapshotsOptionType = azfile.DeleteSnapshotsOptionNone
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

	fileMetadata, err := meta.MustParseJsonMap("file_metadata")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing file_metadata, %w", err)
	}else{
		m.fileMetadata = fileMetadata
	}

	m.rangeSize = int64(meta.ParseInt("range_size", DefaultRangeSize))
	m.parallelism = uint16(meta.ParseInt("parallelism", DefaultParallelism))
	m.count = int64(meta.ParseInt("count", DefaultCount))
	m.offset = int64(meta.ParseInt("offset", DefaultOffset))
	m.maxRetryRequests = meta.ParseInt("max_retry_request", DefaultRetryRequests)
	return m, nil
}
