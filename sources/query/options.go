package query

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/nats-io/nuid"
	"time"
)

const (
	defaultAddress       = "localhost:50000"
	defaultAutoReconnect = true
)

type options struct {
	host                     string
	port                     int
	clientId                 string
	authToken                string
	channel                  string
	group                    string
	autoReconnect            bool
	reconnectIntervalSeconds time.Duration
	maxReconnects            int
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.host, o.port, err = cfg.Properties.MustParseAddress("address", defaultAddress)
	if err != nil {
		return options{}, fmt.Errorf("error parsing address value, %w", err)
	}
	o.authToken = cfg.Properties.ParseString("auth_token", "")

	o.clientId = cfg.Properties.ParseString("client_id", nuid.Next())

	o.channel, err = cfg.Properties.MustParseString("channel")
	if err != nil {
		return o, fmt.Errorf("error parsing channel value, %w", err)
	}

	o.group = cfg.Properties.ParseString("group", "")
	o.autoReconnect = cfg.Properties.ParseBool("auto_reconnect", defaultAutoReconnect)
	interval, err := cfg.Properties.ParseIntWithRange("reconnect_interval_seconds",0, 0, 1000000)
	if err != nil {
		return o, fmt.Errorf("error parsing reconnect interval seconds value, %w", err)
	}
	o.reconnectIntervalSeconds = time.Duration(interval) * time.Second
	o.maxReconnects = cfg.Properties.ParseInt("max_reconnects", 0)
	return o, nil
}
