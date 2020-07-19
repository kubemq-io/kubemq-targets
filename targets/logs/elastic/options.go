package elastic

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"math"
)

type options struct {
	urls                []string
	sniff               bool
	username            string
	password            string
	retryBackoffSeconds int
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}

	var err error
	o.urls, err = cfg.MustParseStringList("urls")
	if err != nil {
		return options{}, fmt.Errorf("error parsing urls, %w", err)
	}
	o.sniff = cfg.ParseBool("sniff", true)
	o.username = cfg.ParseString("username", "")
	o.password = cfg.ParseString("password", "")
	o.retryBackoffSeconds, err = cfg.ParseIntWithRange("retries_backoff_seconds", 0, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing retires backoff seconds, %w", err)
	}

	return o, nil
}
