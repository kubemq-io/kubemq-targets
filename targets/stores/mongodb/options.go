package mongodb

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	url        string
	database   string
	collection string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.url, err = cfg.Properties.MustParseString("url")
	if err != nil {
		return options{}, fmt.Errorf("error parsing url, %w", err)
	}
	o.database, err = cfg.Properties.MustParseString("database")
	if err != nil {
		return options{}, fmt.Errorf("error parsing database, %w", err)
	}
	o.collection, err = cfg.Properties.MustParseString("collection")
	if err != nil {
		return options{}, fmt.Errorf("error parsing collection, %w", err)
	}
	return o, nil
}
