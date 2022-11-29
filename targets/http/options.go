package http

import (
	"fmt"

	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	authType         string
	username         string
	password         string
	token            string
	proxy            string
	rootCertificate  string
	clientPrivateKey string
	clientPublicKey  string
	defaultHeaders   map[string]string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{
		authType:         "",
		username:         "",
		password:         "",
		token:            "",
		proxy:            "",
		rootCertificate:  "",
		clientPublicKey:  "",
		clientPrivateKey: "",
		defaultHeaders:   map[string]string{},
	}
	var err error
	o.authType = cfg.Properties.ParseString("auth_type", "no_auth")
	o.username = cfg.Properties.ParseString("username", "")
	o.password = cfg.Properties.ParseString("password", "")
	o.token = cfg.Properties.ParseString("token", "")
	o.proxy = cfg.Properties.ParseString("proxy", "")
	o.rootCertificate = cfg.Properties.ParseString("root_certificate", "")
	o.clientPrivateKey = cfg.Properties.ParseString("client_private_key", "")
	o.clientPublicKey = cfg.Properties.ParseString("client_public_key", "")
	o.defaultHeaders, err = cfg.Properties.MustParseJsonMap("default_headers")
	if err != nil {
		return options{}, fmt.Errorf("error parsing default_headers value, %w", err)
	}
	return o, nil
}
