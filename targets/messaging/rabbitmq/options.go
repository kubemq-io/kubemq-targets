package rabbitmq

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	url                string
	defaultExchange    string
	defaultTopic       string
	defaultPersistence bool
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.url, err = cfg.Properties.MustParseString("url")
	if err != nil {
		return options{}, fmt.Errorf("error parsing url, %w", err)
	}
	o.defaultExchange = cfg.Properties.ParseString("default_exchange", "")
	o.defaultTopic = cfg.Properties.ParseString("default_topic", "")
	o.defaultPersistence = cfg.Properties.ParseBool("default_persistence", true)

	return o, nil
}

func (o options) defaultMetadata() (metadata, bool) {
	if o.defaultTopic != "" || o.defaultExchange != "" {
		delMode := 2
		if !o.defaultPersistence {
			delMode = 1
		}
		return metadata{
			queue:         o.defaultTopic,
			exchange:      o.defaultExchange,
			mandatory:     false,
			immediate:     false,
			deliveryMode:  delMode,
			priority:      0,
			correlationId: "",
			replyTo:       "",
			expiration:    0,
		}, true
	}
	return metadata{}, false

}
