package spanner

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type metadata struct {
	query     string
	tableName string
	method    string
}

var methodsMap = map[string]string{
	"query":               "query",
	"read":                "read",
	"update_database_ddl": "update_database_ddl",
	"insert":              "insert",
	"update":              "update",
	"insert_or_update":    "insert_or_update",
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}
	if m.method == "query" {
		m.query, err = meta.MustParseString("query")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing query, %w", err)
		}
	}
	if m.method == "read" {
		m.tableName, err = meta.MustParseString("table_name")
		if err != nil {
			return metadata{}, fmt.Errorf("table_name is required for method :%s , error parsing table_name, %w", m.method, err)
		}
	}

	return m, nil
}
