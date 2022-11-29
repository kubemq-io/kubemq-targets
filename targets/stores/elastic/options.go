package elastic

import (
	"fmt"

	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	urls     []string
	sniff    bool
	username string
	password string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}

	var err error
	o.urls, err = cfg.Properties.MustParseStringList("urls")
	if err != nil {
		return options{}, fmt.Errorf("error parsing urls, %w", err)
	}
	o.sniff = cfg.Properties.ParseBool("sniff", true)
	o.username = cfg.Properties.ParseString("username", "")
	o.password = cfg.Properties.ParseString("password", "")

	return o, nil
}
