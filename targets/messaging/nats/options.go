package nats

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/config"
)


const (
	defaultUsername = ""
	defaultPassword = ""
	defaultToken    = ""
	defaultUseTLS   = false
	defaultSSL      = ""
	defaultTimeout  = 100
)


type options struct {
	url            string
	username       string
	password       string
	token          string
	tls            bool
	certFile       string
	certKey        string
	timeout        int
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.url, err = cfg.Properties.MustParseString("url")
	if err != nil {
		return options{}, fmt.Errorf("error parsing url, %w", err)
	}
	o.username = cfg.Properties.ParseString("username", defaultUsername)
	o.password = cfg.Properties.ParseString("password", defaultPassword)
	o.token = cfg.Properties.ParseString("token", defaultToken)

	o.tls = cfg.Properties.ParseBool("tls", defaultUseTLS)
	o.certFile = cfg.Properties.ParseString("cert_file", defaultSSL)
	o.certKey = cfg.Properties.ParseString("cert_key", defaultSSL)

	o.timeout = cfg.Properties.ParseInt("timeout", defaultTimeout)
	if err != nil {
		return options{}, fmt.Errorf("error parsing timeout , %w", err)
	}

	return o, nil
}
