package mysql

import (
	"fmt"
	"math"

	"github.com/kubemq-hub/kubemq-targets/config"
)

const (
	defaultMaxIdleConnections           = 10
	defaultMaxOpenConnections           = 100
	defaultConnectionMaxLifetimeSeconds = 3600
	defaultToken                        = ""
	defaultDBPort = 3306
)

type options struct {
	awsKey                 string
	awsSecretKey           string
	region                 string
	token                  string

	dbPort   int
	dbName   string
	dbUser   string
	endPoint string

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
	o.awsKey, err = cfg.MustParseString("aws_key")
	if err != nil {
		return options{}, fmt.Errorf("error parsing aws_key , %w", err)
	}

	o.awsSecretKey, err = cfg.MustParseString("aws_secret_key")
	if err != nil {
		return options{}, fmt.Errorf("error parsing aws_secret_key , %w", err)
	}

	o.region, err = cfg.MustParseString("region")
	if err != nil {
		return options{}, fmt.Errorf("error parsing region , %w", err)
	}

	o.token = cfg.ParseString("token", defaultToken)

	o.dbName, err = cfg.MustParseString("db_name")
	if err != nil {
		return options{}, fmt.Errorf("error parsing db_name , %w", err)
	}
	o.dbUser, err = cfg.MustParseString("db_user")
	if err != nil {
		return options{}, fmt.Errorf("error parsing db_user , %w", err)
	}
	o.endPoint, err = cfg.MustParseString("end_point")
	if err != nil {
		return options{}, fmt.Errorf("error parsing end_point , %w", err)
	}

	o.dbPort = cfg.ParseInt("db_port",defaultDBPort)
	if err != nil {
		return options{}, fmt.Errorf("error parsing end_point , %w", err)
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
