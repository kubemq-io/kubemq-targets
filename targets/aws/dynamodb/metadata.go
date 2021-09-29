package dynamodb

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/types"
)

type metadata struct {
	method string

	tableName string
}

var methodsMap = map[string]string{
	"list_tables":  "list_tables",
	"create_table": "create_table",
	"delete_table": "delete_table",
	"insert_item":  "insert_item",
	"get_item":     "get_item",
	"delete_item":  "delete_item",
	"update_item":  "update_item",
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}
	if m.method == "insert_item" || m.method == "delete_table" {
		m.tableName, err = meta.MustParseString("table_name")
		if err != nil {
			return metadata{}, fmt.Errorf("table_name is requeired for method:%s, error parsing table_name, %w", m.method, err)
		}
	}
	return m, nil
}
