package aerospike

import (
	"fmt"
	"math"
	"time"

	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	host     string
	port     int
	password string
	username string
	timeout  time.Duration
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.host, err = cfg.Properties.MustParseString("host")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}
	o.port, err = cfg.Properties.MustParseInt("port")
	if err != nil {
		return options{}, fmt.Errorf("error parsing port, %w", err)
	}
	o.username = cfg.Properties.ParseString("username", "")
	o.password = cfg.Properties.ParseString("password", "")
	timeout, err := cfg.Properties.ParseIntWithRange("timeout", 2, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error operation timeout seconds, %w", err)
	}
	o.timeout = time.Duration(timeout) * time.Second
	return o, nil
}
