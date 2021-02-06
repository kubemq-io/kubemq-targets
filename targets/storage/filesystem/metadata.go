package filesystem

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

var methodsMap = map[string]string{
	"save":   "save",
	"load":   "load",
	"delete": "delete",
	"list":   "list",
}

type metadata struct {
	method   string
	path     string
	filename string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing method, %w", err)
	}
	m.filename, err = meta.MustParseString("filename")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing filename, %w", err)
	}
	m.path = meta.ParseString("path", "")
	return m, nil
}
