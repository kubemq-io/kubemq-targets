package postgres

import (
	"database/sql"
	"fmt"

	"github.com/kubemq-io/kubemq-targets/types"
)

var methodsMap = map[string]string{
	"query":       "query",
	"exec":        "exec",
	"transaction": "transaction",
}

var isolationLevelsMap = map[string]string{
	"read_uncommitted": "ReadUncommitted",
	"read_committed":   "ReadCommitted",
	"repeatable_read":  "RepeatableRead",
	"serializable":     "Serializable",
	"":                 "Default",
}

type metadata struct {
	method         string
	isolationLevel sql.IsolationLevel
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}
	isolationLevel, err := meta.ParseStringMap("isolation_level", isolationLevelsMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing isolation_level, %w", err)
	}
	m.isolationLevel = covertToSqlIsolationLevel(isolationLevel)
	return m, nil
}

func covertToSqlIsolationLevel(value string) sql.IsolationLevel {
	switch value {
	case "ReadUncommitted":
		return sql.LevelReadCommitted
	case "ReadCommitted":
		return sql.LevelReadCommitted
	case "RepeatableRead":
		return sql.LevelRepeatableRead
	case "Serializable":
		return sql.LevelSerializable
	default:
		return sql.LevelDefault
	}
}
