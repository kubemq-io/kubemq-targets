package storage

import (
	"github.com/kubemq-hub/kubemq-target-connectors/config"
)

type options struct {
	credentials string
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	var err error
	o.credentials, err = cfg.MustParseString("credentials")
	if err != nil {
		return options{}, err
	}
	return o, nil
}
