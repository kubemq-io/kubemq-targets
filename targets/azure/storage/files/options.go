package blob

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
)

type options struct {
	storageAccessKey string
	storageAccount   string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.storageAccessKey, err = cfg.MustParseString("storage_access_key")
	if err != nil {
		return options{}, fmt.Errorf("error parsing storage_access_key , %w", err)
	}
	o.storageAccount, err = cfg.MustParseString("storage_account")
	if err != nil {
		return options{}, fmt.Errorf("error parsing storage_account , %w", err)
	}
	return o, nil
}
