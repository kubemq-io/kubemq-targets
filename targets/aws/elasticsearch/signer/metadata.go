package signer

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

const (
	DefaultJson = ""
)

type metadata struct {
	method string

	region     string
	json       string

	domain     string
	index      string
	endpoint   string
	service    string
	id         string
}



var httpMethodsMap = map[string]string{
	"GET":     "GET",
	"POST":    "POST",
	"PUT":     "PUT",
	"DELETE":  "DELETE",
	"OPTIONS": "OPTIONS",
}



func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", httpMethodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(httpMethodsMap)
	}
	m.region, err = meta.MustParseString("region")
	if err != nil {
		return metadata{}, fmt.Errorf("error failed to parse region , %w", err)
	}
	if m.method == "GET" {
		m.json = meta.ParseString("json",DefaultJson )
	}else{
		m.json, err = meta.MustParseString("json")
		if err != nil {
			return metadata{}, fmt.Errorf("error failed to parse json , %w", err)
		}
	}
	m.domain, err = meta.MustParseString("domain")
	if err != nil {
		return metadata{}, fmt.Errorf("error failed to parse domain , %w", err)
	}
	m.index, err = meta.MustParseString("index")
	if err != nil {
		return metadata{}, fmt.Errorf("error failed to parse index , %w", err)
	}
	m.endpoint, err = meta.MustParseString("endpoint")
	if err != nil {
		return metadata{}, fmt.Errorf("error failed to parse endpoint , %w", err)
	}
	m.service, err = meta.MustParseString("service")
	if err != nil {
		return metadata{}, fmt.Errorf("error failed to parse service , %w", err)
	}
	m.id, err = meta.MustParseString("id")
	if err != nil {
		return metadata{}, fmt.Errorf("error failed to parse id , %w", err)
	}
	return m, nil
}
