package cassandra

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/kubemq-hub/kubemq-targets/config"
	"math"
)

const (
	defaultProtoVersion      = 4
	defaultReplicationFactor = 1
	defaultConsistency       = gocql.All
	defaultPort              = 9042
)

type options struct {
	hosts             []string
	port              int
	protoVersion      int
	replicationFactor int
	username          string
	password          string
	consistency       gocql.Consistency
	defaultTable      string
	defaultKeyspace   string
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	var err error
	o.hosts, err = cfg.MustParseStringList("hosts")
	if err != nil {
		return options{}, fmt.Errorf("error parsing hosts, %w", err)
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
	o.username = cfg.ParseString("username", "")
	o.password = cfg.ParseString("password", "")
	o.consistency, err = getConsistency(cfg.ParseString("consistency", "All"))
	if err != nil {
		return options{}, fmt.Errorf("error parsing consistency value, %w", err)
	}
	o.defaultTable = cfg.ParseString("default_table", "")
	o.defaultKeyspace = cfg.ParseString("default_keyspace", "")
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
