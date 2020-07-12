package pubsub

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
)

const (
	DefaultRetries = 0
)

type options struct {
	projectID string
	retries   int
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	var err error
	o.projectID, err = cfg.MustParseString("project_id")
	if err != nil {
		return options{}, fmt.Errorf("error parsing project_id, %w", err)
	}
	o.retries = cfg.ParseInt("retries", DefaultRetries)
	err = config.MustExistsEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if err != nil {
		return options{}, err
	}
	return o, nil
}
