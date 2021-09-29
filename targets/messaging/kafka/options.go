package kafka

import (
	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	brokers      []string
	topic        string
	saslUsername string
	saslPassword string
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

	return m, nil
}
