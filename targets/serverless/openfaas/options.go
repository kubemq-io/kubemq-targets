package openfaas

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	gateway  string
	username string
	password string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.gateway, err = cfg.Properties.MustParseString("gateway")
	if err != nil {
		return options{}, fmt.Errorf("error parsing gateway value, %w", err)
	}
	o.username, err = cfg.Properties.MustParseString("username")
	if err != nil {
		return options{}, fmt.Errorf("error parsing username value, %w", err)
	}
	o.password, err = cfg.Properties.MustParseString("password")
	if err != nil {
		return options{}, fmt.Errorf("error parsing password value, %w", err)
	}
	return o, nil
}
