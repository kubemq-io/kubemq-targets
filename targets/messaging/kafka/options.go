package kafka

import (
	"github.com/kubemq-hub/kubemq-target-connectors/config"
)

type options struct {
	brokers      []string
	topic        string
	saslUsername string
	saslPassword string
}

func parseOptions(cfg config.Metadata) (options, error) {
	m := options{}
	var err error
	m.brokers, err = cfg.MustParseStringList("brokers")
	if err != nil {
		return m, err
	}
	m.topic, err = cfg.MustParseString("topic")
	if err != nil {
		return m, err
	}
	m.saslUsername = cfg.ParseString("saslUsername", "")
	m.saslPassword = cfg.ParseString("saslPassword", "")

	return m, nil
}
