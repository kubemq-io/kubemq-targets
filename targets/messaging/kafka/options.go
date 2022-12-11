package kafka

import (
	"strings"

	kafka "github.com/Shopify/sarama"
	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	brokers          []string
	topic            string
	saslUsername     string
	saslPassword     string
	saslMechanism    string
	securityProtocol string
	cacert           string
	clientCert       string
	clientKey        string
	insecure         bool
}

func parseOptions(cfg config.Spec) (options, error) {
	m := options{}
	var err error
	m.brokers, err = cfg.Properties.MustParseStringList("brokers")
	if err != nil {
		return m, err
	}
	m.topic, err = cfg.Properties.MustParseString("topic")
	if err != nil {
		return m, err
	}
	m.saslUsername = cfg.Properties.ParseString("sasl_username", "")
	m.saslPassword = cfg.Properties.ParseString("sasl_password", "")
	m.saslMechanism = cfg.Properties.ParseString("sasl_mechanism", "")
	m.securityProtocol = cfg.Properties.ParseString("security_protocol", "")
	m.cacert = cfg.Properties.ParseString("ca_cert", "")
	m.clientCert = cfg.Properties.ParseString("client_certificate", "")
	m.clientKey = cfg.Properties.ParseString("client_key", "")
	m.insecure = cfg.Properties.ParseBool("insecure", false)

	return m, nil
}

func (m *options) parseASLMechanism() kafka.SASLMechanism {
	switch strings.ToLower(m.saslMechanism) {
	case "plain":
		return kafka.SASLTypePlaintext
	case "scram-sha-256":
		return kafka.SASLTypeSCRAMSHA256
	case "scram-sha-512":
		return kafka.SASLTypeSCRAMSHA512
	case "gssapi", "gss-api", "gss_api":
		return kafka.SASLTypeGSSAPI
	case "oauth", "0auth bearer":
		return kafka.SASLTypeOAuth
	case "external", "ext":
		return kafka.SASLExtKeyAuth
	default:
		return kafka.SASLTypePlaintext
	}
}

func (m *options) parseSecurityProtocol() (bool, bool) {
	switch strings.ToLower(m.securityProtocol) {
	case "plaintext":
		return false, false
	case "ssl":
		return true, false
	case "sasl_plaintext":
		return false, true
	case "sasl_ssl":
		return true, true
	default:
		return false, false
	}
}
