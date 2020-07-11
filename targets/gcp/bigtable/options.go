package bigtable

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
)

type options struct {
	projectID string
	instance string
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	var err error
	o.projectID , err = cfg.MustParseString("project_id")
	if err != nil {
		return options{}, fmt.Errorf("error parsing project_id, %w", err)
	}
	o.instance , err = cfg.MustParseString("instance")
	if err != nil {
		return options{}, fmt.Errorf("error parsing instance, %w", err)
	}
	err = config.MustExistsEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if err !=nil{
		return options{}, err
	}
	return o, nil
}
