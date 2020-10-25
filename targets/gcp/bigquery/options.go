package bigquery

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
)

type options struct {
	projectID   string
	credentials string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.projectID, err = cfg.Properties.MustParseString("project_id")
	if err != nil {
		return options{}, fmt.Errorf("error parsing project_id, %w", err)
	}
	o.credentials, err = cfg.Properties.MustParseString("credentials")
	if err != nil {
		return options{}, fmt.Errorf("error parsing credentials, %w", err)
	}
	return o, nil
}
