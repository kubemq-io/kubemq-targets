package hazelcast

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/types"
)

const (
	defaultListName = ""
	defaultKeyName  = ""
)

var methodsMap = map[string]string{
	"get":      "get",
	"set":      "set",
	"get_list": "get_list",
	"delete":   "delete",
}

type metadata struct {
	method   string
	mapName  string
	key      string
	listName string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing method, %w", err)
	}

	m.key = meta.ParseString("key", defaultKeyName)

	m.mapName, err = meta.MustParseString("map_name")
	if err != nil {
		return metadata{}, fmt.Errorf("error on parsing map_name value, %w", err)
	}

	m.listName = meta.ParseString("list_name", defaultListName)

	return m, nil
}
