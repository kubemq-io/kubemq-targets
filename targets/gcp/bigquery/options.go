package bigquery

import (
	"fmt"

	"github.com/kubemq-io/kubemq-targets/config"
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
	// credentials is optional: when empty the client uses Application Default
	// Credentials (Workload Identity on GKE, the attached service account on
	// Cloud Run / GCE) instead of an inline service-account key.
	o.credentials = cfg.Properties.ParseString("credentials", "")
	return o, nil
}
