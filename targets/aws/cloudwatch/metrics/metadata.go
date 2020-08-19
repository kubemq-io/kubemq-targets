package metrics

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
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
		return metadata{}, meta.GetValidHttpMethodTypes(methodsMap)
	}
	if m.method == "put_metrics" {
		m.namespace, err = meta.MustParseString("namespace")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing namespace, %w", err)
		}
	} else if m.method == "list_metrics" {
		m.namespace = meta.ParseString("namespace", defaultNameSpace)
	}
	return m, nil
}
