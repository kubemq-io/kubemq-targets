package hazelcast

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"time"
)

const (
	defaultUsername                = ""
	defaultPassword                = ""
	defaultUseSSL                  = false
	defaultServerName              = ""
	defaultConnectionAttemptLimit  = 1
	defaultConnectionAttemptPeriod = 36000
	defaultConnectionTimeout       = 36000
)

type options struct {
	address                 string
	username                string
	password                string
	connectionAttemptLimit  int32
	connectionAttemptPeriod time.Duration
	connectionTimeout       time.Duration
	ssl                     bool
	sslcertificatefile      string
	sslcertificatekey       string
	serverName              string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.address, err = cfg.Properties.MustParseString("address")
	if err != nil {
		return options{}, fmt.Errorf("error parsing address, %w", err)
	}
	o.username = cfg.Properties.ParseString("username", defaultUsername)
	o.password = cfg.Properties.ParseString("password", defaultPassword)
	connectionAttemptLimit, err := cfg.Properties.ParseIntWithRange("connection_attempt_limit", defaultConnectionAttemptLimit, 0, 2147483647)
	if err != nil {
		return options{}, fmt.Errorf("error parsing connection_attempt_limit, %w", err)
	}
	o.connectionAttemptLimit = int32(connectionAttemptLimit)

	connectionAttemptPeriod, err := cfg.Properties.ParseIntWithRange("connection_attempt_period", defaultConnectionAttemptPeriod, 0, 2147483647)
	if err != nil {
		return options{}, fmt.Errorf("error parsing connection_attempt_period, %w", err)
	}
	o.connectionAttemptPeriod = time.Duration(connectionAttemptPeriod) * time.Millisecond

	connectionTimeout, err := cfg.Properties.ParseIntWithRange("connection_timeout", defaultConnectionTimeout, 0, 2147483647)
	if err != nil {
		return options{}, fmt.Errorf("error parsing connection_timeout, %w", err)
	}
	o.connectionTimeout = time.Duration(connectionTimeout) * time.Millisecond
	o.ssl = cfg.Properties.ParseBool("ssl", defaultUseSSL)
	o.sslcertificatefile = cfg.Properties.ParseString("cert_file", "")
	o.sslcertificatekey = cfg.Properties.ParseString("cert_key", "")

	o.serverName = cfg.Properties.ParseString("server_name", defaultServerName)

	return o, nil
}
