package redis

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
)

type options struct {
	host      string
	password  string
	enableTLS bool
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{
		host:      "",
		password:  "",
		enableTLS: false,
	}
	var err error
	o.host, err = cfg.MustParseString("host")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}
	o.password = cfg.ParseString("password", "")
	o.enableTLS = cfg.ParseBool("enable_tls", false)
	return o, nil
}
