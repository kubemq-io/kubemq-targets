package lambda

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
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
		i := meta.ParseInt("memory_size", 256)
		m.memorySize = int64(i)
		
		i = meta.ParseInt("timeout", 15)
		m.timeout = int64(i)
		
		m.description = meta.ParseString("description", "")
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
