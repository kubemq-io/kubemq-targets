package elastic

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
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
	o.urls, err = cfg.MustParseStringList("urls")
	if err != nil {
		return options{}, fmt.Errorf("error parsing urls, %w", err)
	}
	o.sniff = cfg.ParseBool("sniff", true)
	o.username = cfg.ParseString("username", "")
	o.password = cfg.ParseString("password", "")

	return o, nil
}
