package rethinkdb

import (
	"fmt"
	"math"
	"time"

	"github.com/kubemq-io/kubemq-targets/config"
)

const (
	defaultUsername         = ""
	defaultPassword         = ""
	defaultAuthKey          = ""
	defaultUseSSL           = false
	defaultSSL              = ""
	defaultTimeout          = 5
	defaultHandShakeVersion = 0
)

type options struct {
	host             string
	username         string
	password         string
	timeout          time.Duration
	keepAlivePeriod  time.Duration
	authKey          string
	ssl              bool
	certFile         string
	certKey          string
	handShakeVersion int
	numberOfRetries  int
	initialCap       int
	maxOpen          int
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.host, err = cfg.Properties.MustParseString("host")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}
	o.username = cfg.Properties.ParseString("username", defaultUsername)
	o.password = cfg.Properties.ParseString("password", defaultPassword)
	timeout, err := cfg.Properties.ParseIntWithRange("timeout", defaultTimeout, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing operation timeout, %w", err)
	}
	o.timeout = time.Duration(timeout) * time.Second

	keepAlivePeriod, err := cfg.Properties.ParseIntWithRange("keep_alive_period", defaultTimeout, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing keep_alive_period, %w", err)
	}
	o.keepAlivePeriod = time.Duration(keepAlivePeriod) * time.Second
	o.authKey = cfg.Properties.ParseString("auth_key", defaultAuthKey)
	o.ssl = cfg.Properties.ParseBool("ssl", defaultUseSSL)
	o.certFile = cfg.Properties.ParseString("cert_file", defaultSSL)
	o.certKey = cfg.Properties.ParseString("cert_key", defaultSSL)
	o.numberOfRetries = cfg.Properties.ParseInt("number_of_retries", 0)
	o.initialCap = cfg.Properties.ParseInt("initial_cap", 0)
	o.maxOpen = cfg.Properties.ParseInt("max_open", 0)
	o.handShakeVersion, err = cfg.Properties.ParseIntWithRange("hand_shake_version", defaultHandShakeVersion, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing hand_shake_version, %w", err)
	}
	return o, nil
}
