package storage

import (
	"github.com/kubemq-hub/kubemq-targets/config"
)

type options struct {
	credentials string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.credentials, err = cfg.Properties.MustParseString("credentials")
	if err != nil {
		return options{}, err
	}
	return o, nil
}
