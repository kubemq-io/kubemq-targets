package redis

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"math"
)

type options struct {
	host                   string
	password               string
	enableTLS              bool
	maxRetries             int
	maxRetryBackoffSeconds int
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{
		host:                   "",
		password:               "",
		enableTLS:              false,
		maxRetries:             0,
		maxRetryBackoffSeconds: 0,
	}
	var err error
	o.host, err = cfg.MustParseString("host")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}
	o.password = cfg.ParseString("password", "")
	o.enableTLS = cfg.ParseBool("enable_tls", false)
	o.maxRetries, err = cfg.ParseIntWithRange("max_retries", 0, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing max retires, %w", err)
	}
	o.maxRetryBackoffSeconds, err = cfg.ParseIntWithRange("max_retries_backoff_seconds", 0, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing max retires backoff seconds, %w", err)
	}
	return o, nil
}
