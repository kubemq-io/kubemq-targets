package cassandra

import (
	"fmt"
	"math"
	"time"

	"github.com/gocql/gocql"
	"github.com/kubemq-io/kubemq-targets/config"
)

const (
	defaultProtoVersion      = 4
	defaultReplicationFactor = 1
	defaultConsistency       = gocql.All
	defaultPort              = 9042
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
	o.username = cfg.Properties.ParseString("username", "")
	o.password = cfg.Properties.ParseString("password", "")
	o.consistency, err = getConsistency(cfg.Properties.ParseString("consistency", "All"))
	if err != nil {
		return options{}, fmt.Errorf("error parsing consistency value, %w", err)
	}
	o.defaultTable = cfg.Properties.ParseString("default_table", "")
	o.defaultKeyspace = cfg.Properties.ParseString("default_keyspace", "")
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
	return o, nil
}

func getConsistency(consistency string) (gocql.Consistency, error) {
	switch consistency {
	case "All":
		return gocql.All, nil
	case "One":
		return gocql.One, nil
	case "Two":
		return gocql.Two, nil
	case "Three":
		return gocql.Three, nil
	case "Quorum":
		return gocql.Quorum, nil
	case "LocalQuorum":
		return gocql.LocalQuorum, nil
	case "EachQuorum":
		return gocql.EachQuorum, nil
	case "LocalOne":
		return gocql.LocalOne, nil
	case "Any":
		return gocql.Any, nil
	case "":
		return defaultConsistency, nil
	}
	return 0, fmt.Errorf("consistency mode %s not found", consistency)
}
