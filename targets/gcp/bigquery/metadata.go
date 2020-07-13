package bigquery

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
)

type metadata struct {
	query     string
	tableName string
	method    string
	datasetID string
}

var methodsMap = map[string]string{
	"query":          "query",
	"create_table":   "create_table",
	"get_table_info": "get_table_info",
	"get_data_sets":  "get_data_sets",
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
	if m.method == "create_table" || m.method == "get_table_info" {
		m.tableName, err = meta.MustParseString("table_name")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing table_name, %w", err)
		}
		m.datasetID, err = meta.MustParseString("dataset_id")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing dataset_id, %w", err)
		}
	}

	return m, nil
}
