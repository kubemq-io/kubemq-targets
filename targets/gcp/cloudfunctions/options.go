package cloudfunctions

import (
	"fmt"

	"github.com/kubemq-hub/kubemq-targets/config"
)

type options struct {
	parrentProject string
	locationMatch  bool
	credentials    string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.parrentProject, err = cfg.MustParseString("project")
	if err != nil {
		return options{}, fmt.Errorf("error parsing project, %w", err)
	}
	o.credentials, err = cfg.MustParseString("credentials")
	if err != nil {
		return options{}, err
	}
	o.locationMatch = cfg.ParseBool("locationMatch", true)

	return o, nil
}
