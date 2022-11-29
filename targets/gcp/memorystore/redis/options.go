package redis

import (
	"fmt"

	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	url string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{
		url: "",
	}
	var err error
	o.url, err = cfg.Properties.MustParseString("url")
	if err != nil {
		return options{}, fmt.Errorf("error parsing url, %w", err)
	}
	return o, nil
}
