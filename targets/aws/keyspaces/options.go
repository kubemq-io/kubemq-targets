package keyspaces

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/kubemq-hub/kubemq-targets/config"
	"math"
	"time"
)

const (
	defaultProtoVersion      = 4
	defaultReplicationFactor = 1
	defaultConsistency       = gocql.LocalQuorum
	defaultPort              = 9142
	defaultUsername          = ""
	defaultPassword          = ""
	defaultKeyspace          = ""
	defaultTable             = ""
)

type options struct {
	hosts                 []string
	port                  int
	protoVersion          int
	replicationFactor     int
	username              string
	password              string
	consistency           gocql.Consistency
	defaultTable          string
	defaultKeyspace       string
	tls                   string
	timeoutSeconds        time.Duration
	connectTimeoutSeconds time.Duration
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.hosts, err = cfg.MustParseStringList("hosts")
	if err != nil {
		return options{}, fmt.Errorf("error parsing hosts, %w", err)
	}
	o.tls, err = cfg.MustParseString("tls")
	if err != nil {
		return options{}, fmt.Errorf("error parsing tls, %w", err)
	}
	o.port, err = cfg.ParseIntWithRange("port", defaultPort, 1, 65535)
	if err != nil {
		return options{}, fmt.Errorf("error parsing port value, %w", err)
	}
	o.protoVersion, err = cfg.ParseIntWithRange("proto_version", defaultProtoVersion, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing proto version value, %w", err)
	}
	o.replicationFactor, err = cfg.ParseIntWithRange("replication_factor", defaultReplicationFactor, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing replication factor value, %w", err)
	}
	o.username = cfg.ParseString("username", defaultUsername)
	o.password = cfg.ParseString("password", defaultPassword)
	o.consistency, err = getConsistency(cfg.ParseString("consistency", "local_quorum"))
	if err != nil {
		return options{}, fmt.Errorf("error parsing consistency value, %w", err)
	}
	o.defaultTable = cfg.ParseString("default_table", defaultTable)
	o.defaultKeyspace = cfg.ParseString("default_keyspace", defaultKeyspace)
	connectTimeout, err := cfg.ParseIntWithRange("connect_timeout_seconds", 60, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing connect timeout seconds timeout value, %w", err)
	}
	o.connectTimeoutSeconds = time.Duration(connectTimeout) * time.Second
	timeout, err := cfg.ParseIntWithRange("timeout_seconds", 60, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing timeout seconds value, %w", err)
	}
	o.timeoutSeconds = time.Duration(timeout) * time.Second
	return o, nil
}

func getConsistency(consistency string) (gocql.Consistency, error) {
	switch consistency {
	case "one":
		return gocql.One, nil
	case "local_quorum":
		return gocql.LocalQuorum, nil
	case "local_one":
		return gocql.LocalOne, nil
	case "":
		return defaultConsistency, nil
	}
	return 0, fmt.Errorf("consistency mode %s not found", consistency)
}
