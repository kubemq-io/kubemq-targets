package memcached

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"math"
)

type options struct {
	hosts                 []string
	maxIdleConnections    int
	defaultTimeoutSeconds int
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{
		hosts:                 nil,
		maxIdleConnections:    0,
		defaultTimeoutSeconds: 0,
	}
	var err error
	o.hosts, err = cfg.MustParseStringList("hosts")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}

	o.maxIdleConnections, err = cfg.ParseIntWithRange("max_idle_connections", 2, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing max idle connections, %w", err)
	}
	o.defaultTimeoutSeconds, err = cfg.ParseIntWithRange("default_timeout_seconds", 1, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing default timeout seconds, %w", err)
	}
	return o, nil
}
