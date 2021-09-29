package consulkv

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/config"
	"time"
)

const (
	defaultUseTLS     = false
	defaultWaitTime   = 36000
	defaultCertFile   = ""
	defaultCertKey    = ""
	defaultScheme     = ""
	defaultDataCenter = ""
	defaultToken      = ""
	defaultTokenFile  = ""
)

type options struct {
	address            string
	scheme             string
	datacenter         string
	token              string
	tokenFile          string
	tls                bool
	tlscertificatefile string
	tlscertificatekey  string
	insecureSkipVerify bool
	waitTime           time.Duration
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.address, err = cfg.Properties.MustParseString("address")
	if err != nil {
		return options{}, fmt.Errorf("error parsing address, %w", err)
	}

	waitTime, err := cfg.Properties.ParseIntWithRange("wait_time", defaultWaitTime, 0, 2147483647)
	if err != nil {
		return options{}, fmt.Errorf("error parsing wait_time, %w", err)
	}
	o.waitTime = time.Duration(waitTime) * time.Millisecond
	o.scheme = cfg.Properties.ParseString("scheme", defaultScheme)
	o.datacenter = cfg.Properties.ParseString("data_center", defaultDataCenter)
	o.token = cfg.Properties.ParseString("token", defaultToken)
	o.tokenFile = cfg.Properties.ParseString("token_file", defaultTokenFile)
	o.insecureSkipVerify = cfg.Properties.ParseBool("insecure_skip_verify", false)


	o.tls = cfg.Properties.ParseBool("tls", defaultUseTLS)
	o.tlscertificatefile = cfg.Properties.ParseString("cert_file", defaultCertFile)
	o.tlscertificatekey = cfg.Properties.ParseString("cert_key", defaultCertKey)
	return o, nil
}
