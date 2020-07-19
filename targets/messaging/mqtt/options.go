package mqtt

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/nats-io/nuid"
)

type options struct {
	host     string
	username string
	password string
	clientId string
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	var err error
	o.host, err = cfg.MustParseString("host")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}
	o.username = cfg.ParseString("username", "")
	o.password = cfg.ParseString("password", "")
	o.clientId = cfg.ParseString("client_id", nuid.Next())
	return o, nil
}
