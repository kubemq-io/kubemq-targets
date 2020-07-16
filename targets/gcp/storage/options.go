package firestore

import (
	"github.com/kubemq-hub/kubemq-target-connectors/config"
)

type options struct {
}

func parseOptions() (options, error) {
	o := options{}
	err := config.MustExistsEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if err != nil {
		return options{}, err
	}
	return o, nil
}
