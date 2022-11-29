package metrics

import (
	"fmt"

	"github.com/kubemq-io/kubemq-targets/types"
)

const (
	defaultNameSpace = ""
)

type metadata struct {
	method    string
	namespace string
}

var methodsMap = map[string]string{
	"put_metrics":  "put_metrics",
	"list_metrics": "list_metrics",
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}
	if m.method == "put_metrics" {
		m.namespace, err = meta.MustParseString("namespace")
		if err != nil {
			return metadata{}, fmt.Errorf("namespace is required for method %s ,error parsing namespace, %w", m.method, err)
		}
	} else if m.method == "list_metrics" {
		m.namespace = meta.ParseString("namespace", defaultNameSpace)
	}
	return m, nil
}
