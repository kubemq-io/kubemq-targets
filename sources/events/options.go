package events

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
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
	responseChannel          string
	autoReconnect            bool
	reconnectIntervalSeconds time.Duration
	maxReconnects            int
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	var err error
	o.host = cfg.ParseString("host", defaultHost)

	o.port, err = cfg.ParseIntWithRange("port", defaultPort, 1, 65535)
	if err != nil {
		return o, fmt.Errorf("error parsing port value, %w", err)
	}

	o.authToken = cfg.ParseString("auth_token", "")

	o.clientId = cfg.ParseString("client_id", nuid.Next())

	o.channel, err = cfg.MustParseString("channel")
	if err != nil {
		return o, fmt.Errorf("error parsing channel value, %w", err)
	}

	o.group = cfg.ParseString("group", "")

	o.concurrency, err = cfg.MustParseIntWithRange("concurrency", 1, 100)
	if err != nil {
		return o, fmt.Errorf("error parsing concurrency value, %w", err)
	}

	o.autoReconnect = cfg.ParseBool("auto_reconnect", defaultAutoReconnect)

	interval, err := cfg.MustParseIntWithRange("reconnect_interval_seconds", 1, 1000000)
	if err != nil {
		return o, fmt.Errorf("error parsing reconnect interval seconds value, %w", err)
	}

	o.reconnectIntervalSeconds = time.Duration(interval) * time.Second

	o.maxReconnects = cfg.ParseInt("max_reconnects", 0)

	o.responseChannel = cfg.ParseString("response_channel", "")

	return o, nil
}
