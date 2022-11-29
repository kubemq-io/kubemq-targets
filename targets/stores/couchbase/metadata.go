package couchbase

import (
	"fmt"
	"math"
	"time"

	"github.com/kubemq-io/kubemq-targets/types"
)

const (
	defaultCas           = 0
	defaultExpirySeconds = 0
)

var methodsMap = map[string]string{
	"get":    "get",
	"set":    "set",
	"delete": "delete",
}

type metadata struct {
	method string
	key    string
	cas    int
	expiry time.Duration
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing method, %w", err)
	}

	m.key, err = meta.MustParseString("key")
	if err != nil {
		return metadata{}, fmt.Errorf("error on parsing key value, %w", err)
	}
	m.cas, err = meta.ParseIntWithRange("cas", defaultCas, 0, math.MaxInt32)
	if err != nil {
		return metadata{}, fmt.Errorf("error on cas value, %w", err)
	}
	expirySeconds, err := meta.ParseIntWithRange("expiry_seconds", defaultExpirySeconds, 0, math.MaxInt32)
	if err != nil {
		return metadata{}, fmt.Errorf("error on expiry seconds value, %w", err)
	}
	m.expiry = time.Duration(expirySeconds) * time.Second
	return m, nil
}
