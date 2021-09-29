package eventhubs

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/config"
)

type options struct {
	connectionString string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	endPoint, err := cfg.Properties.MustParseString("end_point")
	if err != nil {
		return options{}, fmt.Errorf("error parsing end_point , %w", err)
	}
	sharedAccessKeyName, err := cfg.Properties.MustParseString("shared_access_key_name")
	if err != nil {
		return options{}, fmt.Errorf("error parsing shared_access_key_name , %w", err)
	}
	sharedAccessKey, err := cfg.Properties.MustParseString("shared_access_key")
	if err != nil {
		return options{}, fmt.Errorf("error parsing shared_access_key , %w", err)
	}
	entityPath, err := cfg.Properties.MustParseString("entity_path")
	if err != nil {
		return options{}, fmt.Errorf("error parsing entity_path , %w", err)
	}
	o.connectionString = fmt.Sprintf("Endpoint=%s;SharedAccessKeyName=%s;SharedAccessKey=%s;EntityPath=%s", endPoint, sharedAccessKeyName, sharedAccessKey, entityPath)
	return o, nil
}
