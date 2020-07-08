package rabbitmq

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
)

type options struct {
	url string
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	var err error
	o.url, err = cfg.MustParseString("url")
	if err != nil {
		return options{}, fmt.Errorf("error parsing url, %w", err)
	}

	return o, nil
}
