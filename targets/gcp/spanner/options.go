package spanner

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	db          string
	credentials string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.db, err = cfg.Properties.MustParseString("db")
	if err != nil {
		return options{}, fmt.Errorf("error parsing db, %w", err)
	}
	o.credentials, err = cfg.Properties.MustParseString("credentials")
	if err != nil {
		return options{}, fmt.Errorf("error parsing credentials, %w", err)
	}
	return o, nil
}
