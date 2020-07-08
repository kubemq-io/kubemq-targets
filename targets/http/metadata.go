package http

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
	"strings"
)

var methodsMap = map[string]string{
	"post":    "POST",
	"get":     "GET",
	"head":    "HEAD",
	"put":     "PUT",
	"delete":  "DELETE",
	"patch":   "PATCH",
	"options": "OPTIONS",
}

type metadata struct {
	method  string
	url     string
	headers map[string]string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{
		method:  "",
		url:     "",
		headers: map[string]string{},
	}
	var err error
	m.method, err = meta.MustParseString("method")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing method func, %w", err)
	}
	_, ok := methodsMap[strings.ToLower(m.method)]
	if !ok {
		return metadata{}, fmt.Errorf("method %s not supported", m.method)
	}
	m.url, err = meta.MustParseString("url")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing url, %w", err)
	}
	m.headers, err = meta.MustParseJsonMap("headers")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing headers, %w", err)
	}
	return m, nil
}
