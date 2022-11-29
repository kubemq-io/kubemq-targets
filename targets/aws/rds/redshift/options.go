package redshift

import (
	"fmt"
	"math"

	"github.com/kubemq-io/kubemq-targets/config"
)

const (
	defaultMaxIdleConnections           = 10
	defaultMaxOpenConnections           = 100
	defaultConnectionMaxLifetimeSeconds = 3600
)

type options struct {
	connection string
	// maxIdleConnections sets the maximum number of connections in the idle connection pool
	maxIdleConnections int
	// maxOpenConnections sets the maximum number of open connections to the database.
	maxOpenConnections int
	// connectionMaxLifetimeSeconds sets the maximum amount of time a connection may be reused.
	connectionMaxLifetimeSeconds int
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.connection, err = cfg.Properties.MustParseString("connection")
	if err != nil {
		return options{}, fmt.Errorf("error parsing connection string, %w", err)
	}

	o.maxIdleConnections, err = cfg.Properties.ParseIntWithRange("max_idle_connections", defaultMaxIdleConnections, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing max idle connections value, %w", err)
	}
	o.maxOpenConnections, err = cfg.Properties.ParseIntWithRange("max_open_connections", defaultMaxOpenConnections, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing max open connections value, %w", err)
	}
	o.connectionMaxLifetimeSeconds, err = cfg.Properties.ParseIntWithRange("connection_max_lifetime_seconds", defaultConnectionMaxLifetimeSeconds, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing connection max lifetime seconds value, %w", err)
	}
	return o, nil
}
