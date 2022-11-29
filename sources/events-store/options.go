package events_store

import (
	"fmt"
	"time"

	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
)

const (
	defaultAutoReconnect = true
	defaultSources       = 1
)

type options struct {
	host                     string
	port                     int
	clientId                 string
	authToken                string
	channel                  string
	group                    string
	autoReconnect            bool
	responseChannel          string
	reconnectIntervalSeconds time.Duration
	maxReconnects            int
	sources                  int
	doNotParsePayload        bool
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.host, o.port, err = cfg.Properties.MustParseAddress("address", "")
	if err != nil {
		return options{}, fmt.Errorf("error parsing address value, %w", err)
	}
	o.authToken = cfg.Properties.ParseString("auth_token", "")

	o.clientId = cfg.Properties.ParseString("client_id", uuid.New().String())

	o.channel, err = cfg.Properties.MustParseString("channel")
	if err != nil {
		return o, fmt.Errorf("error parsing channel value, %w", err)
	}
	o.sources, err = cfg.Properties.ParseIntWithRange("sources", defaultSources, 1, 1024)
	if err != nil {
		return options{}, fmt.Errorf("error parsing sources value, %w", err)
	}
	o.responseChannel = cfg.Properties.ParseString("response_channel", "")
	o.group = cfg.Properties.ParseString("group", "")
	o.autoReconnect = cfg.Properties.ParseBool("auto_reconnect", defaultAutoReconnect)
	interval, err := cfg.Properties.ParseIntWithRange("reconnect_interval_seconds", 0, 0, 1000000)
	if err != nil {
		return o, fmt.Errorf("error parsing reconnect interval seconds value, %w", err)
	}
	o.reconnectIntervalSeconds = time.Duration(interval) * time.Second
	o.maxReconnects = cfg.Properties.ParseInt("max_reconnects", 0)
	o.doNotParsePayload = cfg.Properties.ParseBool("do_not_parse_payload", false)
	return o, nil
}
