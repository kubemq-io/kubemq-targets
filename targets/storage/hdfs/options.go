package hdfs

import (
	"fmt"

	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	address string
	user    string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.address, err = cfg.Properties.MustParseString("address")
	if err != nil {
		return options{}, fmt.Errorf("error parsing address , %w", err)
	}
	o.user, err = cfg.Properties.MustParseString("user")
	if err != nil {
		return options{}, fmt.Errorf("error parsing user , %w", err)
	}
	return o, nil
}
