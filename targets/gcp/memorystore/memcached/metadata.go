package memcached

import (
	"fmt"

	"github.com/kubemq-io/kubemq-targets/types"
)

var methodsMap = map[string]string{
	"get":    "get",
	"set":    "set",
	"delete": "delete",
}

type metadata struct {
	method string
	key    string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}

	m.key, err = meta.MustParseString("key")
	if err != nil {
		return metadata{}, fmt.Errorf("error on parsing key value, %w", err)
	}

	return m, nil
}
