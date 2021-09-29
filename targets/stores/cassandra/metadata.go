package cassandra

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/types"
)

var methodsMap = map[string]string{
	"get":    "get",
	"set":    "set",
	"delete": "delete",
	"query":  "query",
	"exec":   "exec",
}

var consistencyMap = map[string]string{
	"strong":   "strong",
	"eventual": "eventual",
	"":         "",
}

type metadata struct {
	method      string
	key         string
	consistency string
	table       string
	keyspace    string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing method, %w", err)
	}

	m.key = meta.ParseString("key", "")
	m.consistency, err = meta.ParseStringMap("consistency", consistencyMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error on parsing consistency, %w", err)
	}

	m.table = meta.ParseString("table", "")
	m.keyspace = meta.ParseString("keyspace", "")
	return m, nil
}
func (m metadata) keyspaceTable() string {
	if m.keyspace != "" && m.table != "" {
		return fmt.Sprintf("%s.%s", m.keyspace, m.table)
	}
	return ""
}
