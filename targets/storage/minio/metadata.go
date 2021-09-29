package minio

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/types"
)

var methodsMap = map[string]string{
	"make_bucket":   "make_bucket",
	"list_buckets":  "list_buckets",
	"bucket_exists": "bucket_exists",
	"remove_bucket": "remove_bucket",
	"list_objects":  "list_objects",
	"put":           "put",
	"get":           "get",
	"remove":        "remove",
}

type metadata struct {
	method string
	param1 string
	param2 string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing method, %w", err)
	}
	m.param1 = meta.ParseString("param1", "")
	m.param2 = meta.ParseString("param2", "")
	return m, nil
}
