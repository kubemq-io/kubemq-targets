package http

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/config"
	"math"
)

type options struct {
	authType         string
	username         string
	password         string
	token            string
	proxy            string
	retryCount       int
	retryWaitSeconds int
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
		retryCount:       0,
		retryWaitSeconds: 0,
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
	o.retryCount = cfg.ParseInt("retry_count", -1)
	o.rootCertificate = cfg.ParseString("root_certificate", "")
	o.clientPrivateKey = cfg.ParseString("client_private_key", "")
	o.clientPublicKey = cfg.ParseString("client_public_key", "")
	o.retryWaitSeconds, err = cfg.ParseIntWithRange("retry_wait_seconds", 2, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing retry wait seoncds value, %w", err)
	}
	o.defaultHeaders, err = cfg.MustParseJsonMap("default_headers")
	if err != nil {
		return options{}, fmt.Errorf("error parsing default_headers value, %w", err)
	}
	return o, nil
}
