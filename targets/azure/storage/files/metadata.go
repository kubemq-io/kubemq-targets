package files

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

const (

	DefaultRetryRequests = 1
	DefaultRangeSize     = 4194304
	DefaultFileSize      = 1000000
	DefaultParallelism   = 16

	DefaultCount  = 0
	DefaultOffset = 0
)

var methodsMap = map[string]string{
	"upload": "upload",
	"get":    "get",
	"delete": "delete",
	"create": "create",
}



type metadata struct {
	method                    string
	serviceUrl                string
	rangeSize                 int64
	parallelism               uint16
	offset                    int64
	count                     int64
	size                      int64
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
	m.serviceUrl, err = meta.MustParseString("service_url")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing service_url , %w", err)
	}

	fileMetadata, err := meta.MustParseJsonMap("file_metadata")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing file_metadata, %w", err)
	} else {
		m.fileMetadata = fileMetadata
	}

	m.rangeSize = int64(meta.ParseInt("range_size", DefaultRangeSize))
	m.parallelism = uint16(meta.ParseInt("parallelism", DefaultParallelism))
	m.count = int64(meta.ParseInt("count", DefaultCount))
	m.offset = int64(meta.ParseInt("offset", DefaultOffset))
	m.maxRetryRequests = meta.ParseInt("max_retry_request", DefaultRetryRequests)
	m.size = int64(meta.ParseInt("file_size", DefaultFileSize))
	return m, nil
}
