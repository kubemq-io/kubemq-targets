package bigquery

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type metadata struct {
	query     string
	tableName string
	method    string
	location  string
	datasetID string
}

var methodsMap = map[string]string{
	"query":           "query",
	"create_data_set": "create_data_set",
	"delete_data_set": "delete_data_set",
	"create_table":    "create_table",
	"delete_table":    "delete_table",
	"get_table_info":  "get_table_info",
	"get_data_sets":   "get_data_sets",
	"insert":          "insert",
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
			return metadata{}, fmt.Errorf("query is required for method: %s ,error parsing query, %w", m.method, err)
		}
	}
	if m.method == "create_table" || m.method == "get_table_info" || m.method == "insert" || m.method == "delete_table" {
		m.tableName, err = meta.MustParseString("table_name")
		if err != nil {
			return metadata{}, fmt.Errorf("table_name is required for method: %s ,error parsing table_name, %w", m.method, err)
		}
		m.datasetID, err = meta.MustParseString("dataset_id")
		if err != nil {
			return metadata{}, fmt.Errorf("dataset_id is required for method: %s ,error parsing dataset_id, %w", m.method, err)
		}
	} else if m.method == "create_data_set" || m.method == "delete_data_set" {
		m.datasetID, err = meta.MustParseString("dataset_id")
		if err != nil {
			return metadata{}, fmt.Errorf("dataset_id is required for method: %s ,error parsing dataset_id, %w", m.method, err)
		}
		if m.method == "create_data_set" {
			m.location, err = meta.MustParseString("location")
			if err != nil {
				return metadata{}, fmt.Errorf("location is required for method: %s ,error parsing dataset_id, %w", m.method, err)
			}
		}
	}

	return m, nil
}
