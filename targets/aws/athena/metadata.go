package athena

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type metadata struct {
	method string

	query   string
	catalog string
}

var methodsMap = map[string]string{
	"list_databases":     "list_databases",
	"list_data_catalogs": "list_data_catalogs",
	"query":              "query",
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
	if m.method == "list_databases" {
		m.catalog, err = meta.MustParseString("catalog")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing catalog, %w", err)
		}
	} else if m.method == "query" {
		m.query, err = meta.MustParseString("query")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing query, %w", err)
		}
	}
	return m, nil
}
