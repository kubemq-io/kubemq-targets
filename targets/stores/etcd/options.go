package etcd

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"math"
	"time"
)

type options struct {
	endpoints        []string
	dialTimout       time.Duration
	operationTimeout time.Duration
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	var err error
	o.endpoints, err = cfg.MustParseStringList("endpoints")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}
	dialTimeoutSeconds, err := cfg.ParseIntWithRange("dial_timeout_seconds", 10, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error operation timeout seconds, %w", err)
	}
	o.dialTimout = time.Duration(dialTimeoutSeconds) * time.Second
	operationTimeoutSeconds, err := cfg.ParseIntWithRange("operation_timeout_seconds", 2, 1, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error operation timeout seconds, %w", err)
	}
	o.operationTimeout = time.Duration(operationTimeoutSeconds) * time.Second
	return o, nil
}
