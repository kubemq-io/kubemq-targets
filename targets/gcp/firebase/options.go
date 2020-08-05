package firebase

import (
	"encoding/json"
	"fmt"

	"firebase.google.com/go/v4/messaging"
	"github.com/kubemq-hub/kubemq-targets/config"
)

type options struct {
	projectID        string
	credentials      string
	authClient       bool
	dbClient         bool
	dbURL            string
	messagingClient  bool
	defaultMessaging *messages
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

	o.messagingClient = cfg.ParseBool("messaging_client", false)
	if err != nil {
		return options{}, fmt.Errorf("error parsing messaging_client, %w", err)
	}
	if o.messagingClient {
		o.defaultMessaging = &messages{}
		n := cfg.ParseString("defaultmsg", "")
		if n != "" {
			m := &messaging.Message{}
			err := json.Unmarshal([]byte(n), m)
			if err != nil {
				return o, err
			}
			o.defaultMessaging.single = m
		}

		n = cfg.ParseString("defaultmultimsg", "")
		if n != "" {
			mmulti := &messaging.MulticastMessage{}
			err := json.Unmarshal([]byte(n), mmulti)
			if err != nil {
				return o, err
			}
			o.defaultMessaging.multicast = mmulti
		}
	}
	return o, nil
}
