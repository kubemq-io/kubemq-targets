package messaging

import (
	"encoding/json"
	"fmt"

	"firebase.google.com/go/messaging"
	"github.com/kubemq-hub/kubemq-targets/config"
)

type options struct {
	projectID       string
	credentials     string
	defult          *messaging.Message
	defaultmultimsg *messaging.MulticastMessage
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

	n := cfg.ParseString("defaultmsg", "")
	if n != "" {
		err := json.Unmarshal([]byte(n), &o.defult)
		if err != nil {
			return o, err
		}
	}

	n := cfg.ParseString("defaultmultimsg", "")
	if n != "" {
		err := json.Unmarshal([]byte(n), &o.defaultmultimsg)
		if err != nil {
			return o, err
		}
	}
	return o, nil
}
