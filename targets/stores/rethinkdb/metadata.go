package rethinkdb

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/types"
)

const (
	defaultKey = ""
)

var methodsMap = map[string]string{
	"get":    "get",
	"update": "update",
	"insert": "insert",
	"delete": "delete",
}

type metadata struct {
	method string
	key    string
	dbName string
	table  string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing method, %w", err)
	}
	m.dbName, err = meta.MustParseString("db_name")
	if err != nil {
		return metadata{}, fmt.Errorf("error db_name, %w", err)
	}
	m.table, err = meta.MustParseString("table")
	if err != nil {
		return metadata{}, fmt.Errorf("error table, %w", err)
	}

	m.key = meta.ParseString("key", defaultKey)

	return m, nil
}
