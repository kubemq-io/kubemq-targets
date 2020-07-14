package http

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
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

func parseOptions(cfg config.Metadata) (options, error) {
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
	o.authType = cfg.ParseString("auth_type", "no_auth")
	o.username = cfg.ParseString("username", "")
	o.password = cfg.ParseString("password", "")
	o.token = cfg.ParseString("token", "")
	o.proxy = cfg.ParseString("proxy", "")
	o.rootCertificate = cfg.ParseString("root_certificate", "")
	o.clientPrivateKey = cfg.ParseString("client_private_key", "")
	o.clientPublicKey = cfg.ParseString("client_public_key", "")
	o.defaultHeaders, err = cfg.MustParseJsonMap("default_headers")
	if err != nil {
		return options{}, fmt.Errorf("error parsing default_headers value, %w", err)
	}
	return o, nil
}
