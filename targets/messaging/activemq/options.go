package activemq

import (
	"fmt"

	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	host               string
	username           string
	password           string
	defaultDestination string
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
	o.defaultDestination = cfg.Properties.ParseString("default_destination", "")
	return o, nil
}

func (o options) defaultMetadata() (metadata, bool) {
	if o.defaultDestination != "" {
		return metadata{
			destination: o.defaultDestination,
		}, true
	}
	return metadata{}, false
}
