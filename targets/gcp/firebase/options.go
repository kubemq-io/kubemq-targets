package firebase

import (
	"encoding/json"
	"fmt"

	"firebase.google.com/go/v4/messaging"
	"github.com/kubemq-io/kubemq-targets/config"
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
	o.projectID, err = cfg.Properties.MustParseString("project_id")
	if err != nil {
		return options{}, fmt.Errorf("error parsing project_id, %w", err)
	}
	o.credentials, err = cfg.Properties.MustParseString("credentials")
	if err != nil {
		return options{}, err
	}
	o.authClient, err = cfg.Properties.MustParseBool("auth_client")
	if err != nil {
		return options{}, fmt.Errorf("error parsing auth_client, %w", err)
	}

	o.dbClient, err = cfg.Properties.MustParseBool("db_client")
	if err != nil {
		return options{}, fmt.Errorf("error parsing db_client, %w", err)
	}
	o.dbURL = cfg.Properties.ParseString("db_url", "")

	o.messagingClient, err = cfg.Properties.MustParseBool("messaging_client")
	if err != nil {
		return options{}, fmt.Errorf("error parsing messaging_client, %w", err)
	}
	if o.messagingClient {
		o.defaultMessaging = &messages{}
		n := cfg.Properties.ParseString("defaultmsg", "")
		if n != "" {
			m := &messaging.Message{}
			err := json.Unmarshal([]byte(n), m)
			if err != nil {
				return o, err
			}
			o.defaultMessaging.single = m
		}

		n = cfg.Properties.ParseString("defaultmultimsg", "")
		if n != "" {
			multi := &messaging.MulticastMessage{}
			err := json.Unmarshal([]byte(n), multi)
			if err != nil {
				return o, err
			}
			o.defaultMessaging.multicast = multi
		}
	}
	return o, nil
}
