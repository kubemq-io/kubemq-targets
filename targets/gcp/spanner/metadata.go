package spanner

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
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
	if m.method == "query" {
		m.query, err = meta.MustParseString("query")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing query, %w", err)
		}
	}
	if m.method == "read" {
		m.tableName, err = meta.MustParseString("table_name")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing query, %w", err)
		}
	}

	return m, nil
}
