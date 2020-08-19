package redis

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
	"math"
)

var methodsMap = map[string]string{
	"get":    "get",
	"set":    "set",
	"delete": "delete",
}
var concurrencyMap = map[string]string{
	"first-write": "first-write",
	"last-write":  "last-write",
	"":            "",
}
var consistencyMap = map[string]string{
	"strong":   "strong",
	"eventual": "eventual",
	"":         "",
}

type metadata struct {
	method      string
	key         string
	etag        int
	concurrency string
	consistency string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidHttpMethodTypes(methodsMap)
	}

	m.key, err = meta.MustParseString("key")
	if err != nil {
		return metadata{}, fmt.Errorf("error on parsing key value, %w", err)
	}
	m.etag, err = meta.ParseIntWithRange("etag", 0, 0, math.MaxInt32)
	if err != nil {
		return metadata{}, fmt.Errorf("error on parsing etag value, %w", err)
	}
	m.concurrency, err = meta.ParseStringMap("concurrency", concurrencyMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error on parsing concurrency, %w", err)
	}

	m.consistency, err = meta.ParseStringMap("consistency", consistencyMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error on parsing consistency, %w", err)
	}
	return m, nil
}
