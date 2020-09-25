package elasticsearch

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
)

const (
	DefaultToken = ""
)

type options struct {
	awsKey       string
	awsSecretKey string
	token        string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.awsKey, err = cfg.Properties.MustParseString("aws_key")
	if err != nil {
		return options{}, fmt.Errorf("error parsing aws_key , %w", err)
	}

	o.awsSecretKey, err = cfg.Properties.MustParseString("aws_secret_key")
	if err != nil {
		return options{}, fmt.Errorf("error parsing aws_secret_key , %w", err)
	}

	o.token = cfg.Properties.ParseString("token", DefaultToken)

	return o, nil
}
