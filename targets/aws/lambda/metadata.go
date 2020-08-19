package lambda

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

const (
	defaultDescription           = ""
	defaultMemorySize           = 256
	defaultTimeout = 15
)


type metadata struct {
	method string
	
	zipFileName  string
	functionName string
	handlerName  string
	role         string
	runtime      string
	memorySize   int64
	timeout      int64
	description  string
}

var methodsMap = map[string]string{
	"list":   "list",
	"create": "create",
	"run":    "run",
	"delete": "delete",
}


func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}
	if m.method == "create" {
		m.zipFileName, err = meta.MustParseString("zip_file_name")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing zip_file_name, %w", err)
		}
		m.handlerName, err = meta.MustParseString("handler_name")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing handler_name, %w", err)
		}
		m.role, err = meta.MustParseString("role")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing role, %w", err)
		}
		m.runtime, err = meta.MustParseString("runtime")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing runtime, %w", err)
		}
		i := meta.ParseInt("memory_size", defaultMemorySize)
		m.memorySize = int64(i)
		
		i = meta.ParseInt("timeout", defaultTimeout)
		m.timeout = int64(i)
		
		m.description = meta.ParseString("description", defaultDescription)
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing runtime, %w", err)
		}
	}
	if m.method != "list" {
		m.functionName, err = meta.MustParseString("function_name")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing function_name, %w", err)
		}
	}
	return m, nil
}
