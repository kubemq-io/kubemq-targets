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
	o.hosts, err = cfg.Properties.MustParseStringList("hosts")
	if err != nil {
		return options{}, fmt.Errorf("error parsing hosts, %w", err)
	}
	o.port, err = cfg.Properties.ParseIntWithRange("port", defaultPort, 1, 65535)
	if err != nil {
		return options{}, fmt.Errorf("error parsing port value, %w", err)
	}
	o.protoVersion, err = cfg.Properties.ParseIntWithRange("proto_version", defaultProtoVersion, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing proto version value, %w", err)
	}
	o.replicationFactor, err = cfg.Properties.ParseIntWithRange("replication_factor", defaultReplicationFactor, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing replication factor value, %w", err)
	}
	o.username = cfg.Properties.ParseString("username", defaultUsername)
	o.password = cfg.Properties.ParseString("password", defaultPassword)
	o.consistency, err = getConsistency(cfg.Properties.ParseString("consistency", "LocalQuorum"))
	if err != nil {
		return options{}, fmt.Errorf("error parsing consistency value, %w", err)
	}
	o.defaultTable = cfg.Properties.ParseString("default_table", defaultTable)
	o.defaultKeyspace = cfg.Properties.ParseString("default_keyspace", defaultKeyspace)
	connectTimeout, err := cfg.Properties.ParseIntWithRange("connect_timeout_seconds", 60, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing connect timeout seconds timeout value, %w", err)
	}
	o.connectTimeoutSeconds = time.Duration(connectTimeout) * time.Second
	timeout, err := cfg.Properties.ParseIntWithRange("timeout_seconds", 60, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing timeout seconds value, %w", err)
	}
	o.timeoutSeconds = time.Duration(timeout) * time.Second
	o.tls, err = cfg.Properties.MustParseString("tls")
	if err != nil {
		return options{}, fmt.Errorf("error parsing tls, %w", err)
	}
	return o, nil
}

func getConsistency(consistency string) (gocql.Consistency, error) {
	switch consistency {
	case "One":
		return gocql.One, nil
	case "LocalQuorum":
		return gocql.LocalQuorum, nil
	case "LocalOne":
		return gocql.LocalOne, nil
	case "":
		return defaultConsistency, nil
	}
	return 0, fmt.Errorf("consistency mode %s not found", consistency)
}
