package postgres

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"math"
)

const (
	defaultMaxIdleConnections           = 10
	defaultMaxOpenConnections           = 100
	defaultConnectionMaxLifetimeSeconds = 3600
)

type options struct {
	credentials string
	useProxy    bool
	instanceConnectionName string
	dbUser                 string
	dbName                 string
	dbPassword             string
	connection string
	// maxIdleConnections sets the maximum number of connections in the idle connection pool
	maxIdleConnections int
	//maxOpenConnections sets the maximum number of open connections to the database.
	maxOpenConnections int
	// connectionMaxLifetimeSeconds sets the maximum amount of time a connection may be reused.
	connectionMaxLifetimeSeconds int
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error

	o.useProxy = cfg.ParseBool("use_proxy", true)
	if o.useProxy {
		o.instanceConnectionName, err = cfg.MustParseString("instance_connection_name")
		if err != nil {
			return options{}, fmt.Errorf("error parsing instance_connection_name string, %w", err)
		}
		o.dbUser, err = cfg.MustParseString("db_user")
		if err != nil {
			return options{}, fmt.Errorf("error parsing db_user string, %w", err)
		}
		o.dbName, err = cfg.MustParseString("db_name")
		if err != nil {
			return options{}, fmt.Errorf("error parsing db_name string, %w", err)
		}
		o.dbPassword, err = cfg.MustParseString("db_password")
		if err != nil {
			return options{}, fmt.Errorf("error parsing db_password string, %w", err)
		}
		o.credentials, err = cfg.MustParseString("credentials")
		if err != nil {
			return options{}, err
		}
	} else {
		o.connection, err = cfg.MustParseString("connection")
		if err != nil {
			return options{}, fmt.Errorf("error parsing connection string, %w", err)
		}
	}

	o.maxIdleConnections, err = cfg.ParseIntWithRange("max_idle_connections", defaultMaxIdleConnections, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing max idle connections value, %w", err)
	}
	o.maxOpenConnections, err = cfg.ParseIntWithRange("max_open_connections", defaultMaxOpenConnections, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing max open connections value, %w", err)
	}
	o.connectionMaxLifetimeSeconds, err = cfg.ParseIntWithRange("connection_max_lifetime_seconds", defaultConnectionMaxLifetimeSeconds, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing connection max lifetime seconds value, %w", err)
	}

	return o, nil
}
