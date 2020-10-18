package msk

import (
	"github.com/kubemq-hub/kubemq-targets/config"
)

const (
	DefaultSaslUsername = ""
	DefaultSaslPassword = ""
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
	m.saslUsername = cfg.Properties.ParseString("saslUsername", DefaultSaslUsername)
	m.saslPassword = cfg.Properties.ParseString("saslPassword", DefaultSaslPassword)

	return m, nil
}
