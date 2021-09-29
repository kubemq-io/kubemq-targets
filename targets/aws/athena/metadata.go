package athena

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/types"
)

type metadata struct {
	method string

	query          string
	catalog        string
	outputLocation string
	DB             string
	executionID    string
}

var methodsMap = map[string]string{
	"list_databases":     "list_databases",
	"list_data_catalogs": "list_data_catalogs",
	"query":              "query",
	"get_query_result":   "get_query_result",
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}
	if m.method != "list_data_catalogs" && m.method != "get_query_result" {
		m.catalog, err = meta.MustParseString("catalog")
		if err != nil {
			return metadata{}, fmt.Errorf("catalog is required for method:%s,error parsing catalog, %w", m.method, err)
		}
		if m.method == "query" {
			m.query, err = meta.MustParseString("query")
			if err != nil {
				return metadata{}, fmt.Errorf("query is required for method:%s , error parsing query, %w", m.method, err)
			}
			m.DB, err = meta.MustParseString("db")
			if err != nil {
				return metadata{}, fmt.Errorf("db is required for method:%s , error parsing db, %w", m.method, err)
			}
			m.outputLocation, err = meta.MustParseString("output_location")
			if err != nil {
				return metadata{}, fmt.Errorf("output_location is required for method:%s ,error parsing output_location, %w", m.method, err)
			}
		}
	}
	if m.method == "get_query_result" {
		m.executionID, err = meta.MustParseString("execution_id")
		if err != nil {
			return metadata{}, fmt.Errorf("execution_id is required for method:%s error parsing execution_id, %w", m.method, err)
		}
	}
	return m, nil
}
