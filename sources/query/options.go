package query

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/nats-io/nuid"
	"time"
)

type options struct {
	host                     string
	port                     int
	clientId                 string
	authToken                string
	channel                  string
	group                    string
	concurrency              int
	autoReconnect            bool
	reconnectIntervalSeconds time.Duration
	maxReconnects            int
}

func parseOptions(cfg config.Metadata) (options, error) {
	m := options{}
	var err error
	m.host = cfg.ParseString("host", defaultHost)

	m.port, err = cfg.ParseIntWithRange("port", defaultPort, 1, 65535)
	if err != nil {
		return m, fmt.Errorf("error parsing port value, %w", err)
	}

	m.authToken = cfg.ParseString("auth_token", "")

	m.clientId = cfg.ParseString("client_id", nuid.Next())

	m.channel, err = cfg.MustParseString("channel")
	if err != nil {
		return m, fmt.Errorf("error parsing channel value, %w", err)
	}

	m.group = cfg.ParseString("group", "")

	m.concurrency, err = cfg.MustParseIntWithRange("concurrency", 1, 100)
	if err != nil {
		return m, fmt.Errorf("error parsing concurrency value, %w", err)
	}

	m.autoReconnect = cfg.ParseBool("auto_reconnect", defaultAutoReconnect)

	interval, err := cfg.MustParseIntWithRange("reconnect_interval_seconds", 1, 1000000)
	if err != nil {
		return m, fmt.Errorf("error parsing reconnect interval seconds value, %w", err)
	}

	m.reconnectIntervalSeconds = time.Duration(interval) * time.Second

	m.maxReconnects = cfg.ParseInt("max_reconnects", 0)

	return m, nil
}
