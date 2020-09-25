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
	o.parrentProject, err = cfg.Properties.MustParseString("project")
	if err != nil {
		return options{}, fmt.Errorf("error parsing project, %w", err)
	}
	o.credentials, err = cfg.Properties.MustParseString("credentials")
	if err != nil {
		return options{}, err
	}
	o.locationMatch = cfg.Properties.ParseBool("location_match", true)

	return o, nil
}
