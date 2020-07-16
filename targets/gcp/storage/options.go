package storage

import (
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"os"
)

type options struct {
	credentials string
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	o.credentials = cfg.ParseString("credentials", "")
	if o.credentials != "" {
		_ = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", o.credentials)
	}
	err := config.MustExistsEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if err != nil {
		return options{}, err
	}
	return o, nil
}
