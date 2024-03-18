package amqp

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Azure/go-amqp"
	"strings"

	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	url      string
	username string
	password string
	caCert   string
	insecure bool
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.url, err = cfg.Properties.MustParseString("url")
	if err != nil {
		return options{}, fmt.Errorf("error parsing url, %w", err)
	}
	o.username = cfg.Properties.ParseString("username", "")
	o.password = cfg.Properties.ParseString("password", "")
	o.caCert = cfg.Properties.ParseString("ca_cert", "")
	o.insecure = cfg.Properties.ParseBool("skip_insecure", false)
	return o, nil
}

func (o options) getConnOptions() (*amqp.ConnOptions, error) {
	connOptions := &amqp.ConnOptions{}
	if strings.HasPrefix(o.url, "amqps://") {
		tlsCfg := &tls.Config{
			InsecureSkipVerify: o.insecure,
		}
		if o.caCert != "" {
			caCertPool := x509.NewCertPool()
			if !caCertPool.AppendCertsFromPEM([]byte(o.caCert)) {
				return nil, fmt.Errorf("error loading Root CA Cert")
			}
			tlsCfg.RootCAs = caCertPool
		}
		connOptions.TLSConfig = tlsCfg
	}
	if o.username != "" {
		connOptions.SASLType = amqp.SASLTypePlain(o.username, o.password)
	} else {
		connOptions.SASLType = amqp.SASLTypeAnonymous()
	}
	return connOptions, nil
}
