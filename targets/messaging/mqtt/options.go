package mqtt

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/config"
	"github.com/kubemq-io/kubemq-targets/pkg/uuid"
)

type options struct {
	host         string
	username     string
	password     string
	clientId     string
	defaultTopic string
	defaultQos   int
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.host, err = cfg.Properties.MustParseString("host")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}
	o.username = cfg.Properties.ParseString("username", "")
	o.password = cfg.Properties.ParseString("password", "")
	o.clientId = cfg.Properties.ParseString("client_id", uuid.New().String())
	o.defaultTopic = cfg.Properties.ParseString("default_topic", "")
	o.defaultQos, err = cfg.Properties.ParseIntWithRange("default_qos", 0, 0, 2)
	if err != nil {
		return options{}, fmt.Errorf("error parsing default_qos, %w", err)
	}
	return o, nil
}
func (o options) defaultMetadata() (metadata, bool) {
	if o.defaultTopic != "" {
		return metadata{
			topic: o.defaultTopic,
			qos:   o.defaultQos,
		}, true
	}
	return metadata{}, false

}
