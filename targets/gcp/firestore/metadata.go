package firestore

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

var methodsMap = map[string]string{
	"documents_all":       "documents_all",
	"document_key":        "document_key",
	"delete_document_key": "delete_document_key",
	"add":                 "add",
}

type metadata struct {
	method string
	key    string
	item   string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}

	m.key, err = meta.MustParseString("collection")
	if err != nil {
		return metadata{}, fmt.Errorf("error on parsing collection value, %w", err)
	}
	if m.method == "document_key" || m.method == "delete_document_key" {
			m.item, err = meta.MustParseString("item")
		if err != nil {
			return metadata{}, fmt.Errorf("item is required for method: %s,error on parsing item value, %w", m.method, err)
		}
	}

	return m, nil
}
