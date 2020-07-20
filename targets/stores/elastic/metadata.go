package elastic

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

var methodsMap = map[string]string{
	"get":    "get",
	"set":    "set",
	"delete": "delete",
}

type metadata struct {
	method string
	index  string
	id     string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing method, %w", err)
	}
	m.index, err = meta.MustParseString("index")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing index value, %w", err)
	}
	m.id, err = meta.MustParseString("id")
	if err != nil {
		return metadata{}, fmt.Errorf("error on parsing id value, %w", err)
	}
	return m, nil
}
