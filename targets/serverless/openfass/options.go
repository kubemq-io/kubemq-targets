package openfass

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
)

type options struct {
	gateway  string
	username string
	password string
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	var err error
	o.gateway, err = cfg.MustParseString("gateway")
	if err != nil {
		return options{}, fmt.Errorf("error parsing gateway value, %w", err)
	}
	o.username, err = cfg.MustParseString("username")
	if err != nil {
		return options{}, fmt.Errorf("error parsing username value, %w", err)
	}
	o.password, err = cfg.MustParseString("password")
	if err != nil {
		return options{}, fmt.Errorf("error parsing password value, %w", err)
	}
	return o, nil
}
