package amazonmq

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
)



type options struct {
	host     string
	username string
	password string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.host, err = cfg.MustParseString("host")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}
	o.username ,err= cfg.MustParseString("username")
	if err != nil {
		return options{}, fmt.Errorf("error parsing username , %w", err)
	}
	o.password ,err= cfg.MustParseString("password")
	if err != nil {
		return options{}, fmt.Errorf("error parsing password , %w", err)
	}
	return o, nil
}
