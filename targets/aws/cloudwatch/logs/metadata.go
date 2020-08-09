package logs

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

const (
	defaultLimit = 100
)


type metadata struct {
	method string

	limit       int64
	groupName   string
	streamName  string
	groupPrefix string
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
	return m, nil
}
