package consulkv

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
	"time"
)

const (
	defaultKey               = ""
	defaultNear              = ""
	defaultFilter            = ""
	defaultPrefix            = ""
	defaultAllowStale        = false
	defaultRequireConsistent = false
	defaultUserCache         = false

	defaultMaxAge       = 36000
	defaultStaleIfError = 36000
)

var methodsMap = map[string]string{
	"get":    "get",
	"put":    "put",
	"list":   "list",
	"delete": "delete",
}

type metadata struct {
	method            string
	key               string
	near              string
	filter            string
	prefix            string
	allowStale        bool
	requireConsistent bool
	useCache          bool
	maxAge            time.Duration
	staleIfError      time.Duration
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing method, %w", err)
	}

	m.key = meta.ParseString("key", defaultKey)
	m.near = meta.ParseString("near", defaultNear)
	m.filter = meta.ParseString("filter", defaultFilter)
	m.prefix = meta.ParseString("prefix", defaultPrefix)

	m.allowStale = meta.ParseBool("allow_stale", defaultAllowStale)
	m.requireConsistent = meta.ParseBool("require_consistent", defaultRequireConsistent)
	m.useCache = meta.ParseBool("user_cache", defaultUserCache)


	maxAge, err := meta.ParseIntWithRange("max_age", defaultMaxAge, 0, 2147483647)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing max_age, %w", err)
	}
	m.maxAge = time.Duration(maxAge) * time.Millisecond

	staleIfError, err := meta.ParseIntWithRange("stale_if_error", defaultStaleIfError, 0, 2147483647)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing stale_if_error, %w", err)
	}
	m.staleIfError = time.Duration(staleIfError) * time.Millisecond

	return m, nil
}
