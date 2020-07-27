package firebase

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
)

type options struct {
	projectID   string
	credentials string
	authClient  bool
	dbClient    bool
	dbURL       string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.projectID, err = cfg.MustParseString("project_id")
	if err != nil {
		return options{}, fmt.Errorf("error parsing project_id, %w", err)
	}
	o.credentials, err = cfg.MustParseString("credentials")
	if err != nil {
		return options{}, err
	}
	o.authClient = cfg.ParseBool("auth_client", false)
	if err != nil {
		return options{}, fmt.Errorf("error parsing auth_client, %w", err)
	}

		o.dbClient = cfg.ParseBool("db_client", false)
	if err != nil {
		return options{}, fmt.Errorf("error parsing db_client, %w", err)
	}
	o.dbURL = cfg.ParseString("db_url", "")
	return o, nil
}
