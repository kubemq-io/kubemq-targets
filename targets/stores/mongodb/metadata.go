package mongodb

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/types"
)

var methodsMap = map[string]string{
	"get_by_key":    "get_by_key",
	"set_by_key":    "set_by_key",
	"delete_by_key": "delete_by_key",
	"find":          "find",
	"find_many":     "find_many",
	"insert":        "insert",
	"insert_many":   "insert_many",
	"update":        "update",
	"update_many":   "update_many",
	"delete":        "delete",
	"delete_many":   "delete_many",
	"aggregate":     "aggregate",
	"distinct":      "distinct",
}

type metadata struct {
	method    string
	key       string
	filter    map[string]interface{}
	fieldName string
	setUpsert bool
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing method, %w", err)
	}
	m.key = meta.ParseString("key", "")
	m.fieldName = meta.ParseString("field_name", "")
	m.filter, err = meta.MustParseInterfaceMap("filter")

	if err != nil {
		return metadata{}, fmt.Errorf("error parsing filter, %w", err)
	}
	m.setUpsert = meta.ParseBool("set_upsert", false)
	return m, nil
}
