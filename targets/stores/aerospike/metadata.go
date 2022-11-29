package aerospike

import (
	"fmt"

	"github.com/kubemq-io/kubemq-targets/types"
)

var methodsMap = map[string]string{
	"get":       "get",
	"set":       "set",
	"delete":    "delete",
	"get_batch": "get_batch",
}

type metadata struct {
	method    string
	key       string
	userKey   string
	namespace string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing method, %w", err)
	}

	m.key = meta.ParseString("key", "")
	m.userKey = meta.ParseString("user_key", "")
	m.namespace = meta.ParseString("namespace", "")

	return m, nil
}
