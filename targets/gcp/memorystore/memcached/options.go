package memcached

import (
	"fmt"
	"math"

	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	hosts                 []string
	maxIdleConnections    int
	defaultTimeoutSeconds int
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{
		hosts:                 nil,
		maxIdleConnections:    0,
		defaultTimeoutSeconds: 0,
	}
	var err error
	o.hosts, err = cfg.Properties.MustParseStringList("hosts")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}

	o.maxIdleConnections, err = cfg.Properties.ParseIntWithRange("max_idle_connections", 2, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing max idle connections, %w", err)
	}
	o.defaultTimeoutSeconds, err = cfg.Properties.ParseIntWithRange("default_timeout_seconds", 1, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing default timeout seconds, %w", err)
	}
	return o, nil
}
