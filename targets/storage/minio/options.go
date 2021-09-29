package minio

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	endpoint        string
	useSSL          bool
	accessKeyId     string
	secretAccessKey string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.endpoint, err = cfg.Properties.MustParseString("endpoint")
	if err != nil {
		return options{}, fmt.Errorf("error parsing endpoint, %w", err)
	}
	o.useSSL = cfg.Properties.ParseBool("use_ssl", true)
	o.accessKeyId, err = cfg.Properties.MustParseEnv("access_key_id", "ACCESS_KEY_ID", "")
	if err != nil {
		return options{}, fmt.Errorf("error parsing access key id, %w", err)
	}

	o.secretAccessKey, err = cfg.Properties.MustParseEnv("secret_access_key", "SECRET_ACCESS_KEY", "")
	if err != nil {
		return options{}, fmt.Errorf("error parsing secret access key, %w", err)
	}
	return o, nil
}
