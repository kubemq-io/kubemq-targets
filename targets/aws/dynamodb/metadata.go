package dynamodb

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
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

func getValidMethodTypes() string {
	s := "invalid method type, method type should be one of the following:"
	for k := range methodsMap {
		s = fmt.Sprintf("%s :%s,", s, k)
	}
	return s
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, fmt.Errorf(getValidMethodTypes())
	}
	if m.method == "insert_item" || m.method == "delete_table" {
		m.tableName, err = meta.MustParseString("table_name")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing table_name, %w", err)
		}
	}
	return m, nil
}
