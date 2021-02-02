package queue

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/kubemq-hub/kubemq-targets/pkg/uuid"
)

const (
	defaultAddress     = "localhost:50000"
	defaultBatchSize   = 1
	defaultWaitTimeout = 60
	defaultSources     = 1
)

type options struct {
	host            string
	port            int
	clientId        string
	authToken       string
	channel         string
	sources         int
	responseChannel string
	batchSize       int
	waitTimeout     int
	maxRequeue      int
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.host, o.port, err = cfg.Properties.MustParseAddress("address", defaultAddress)
	if err != nil {
		return options{}, fmt.Errorf("error parsing address value, %w", err)
	}
	o.authToken = cfg.Properties.ParseString("auth_token", "")

	o.clientId = cfg.Properties.ParseString("client_id", uuid.New().String())

	o.channel, err = cfg.Properties.MustParseString("channel")
	if err != nil {
		return options{}, fmt.Errorf("error parsing channel value, %w", err)
	}
	o.sources, err = cfg.Properties.ParseIntWithRange("sources", defaultSources, 1, 1024)
	if err != nil {
		return options{}, fmt.Errorf("error parsing sources value, %w", err)
	}

	o.responseChannel = cfg.Properties.ParseString("response_channel", "")

	o.batchSize, err = cfg.Properties.ParseIntWithRange("batch_size", defaultBatchSize, 1, 1024)
	if err != nil {
		return options{}, fmt.Errorf("error parsing batch size value, %w", err)
	}
	o.waitTimeout, err = cfg.Properties.ParseIntWithRange("wait_timeout", defaultWaitTimeout, 1, 24*60*60)
	if err != nil {
		return options{}, fmt.Errorf("error parsing wait timeout value, %w", err)
	}
	o.maxRequeue, err = cfg.Properties.ParseIntWithRange("max_requeue", 0, 0, 1024)
	if err != nil {
		return options{}, fmt.Errorf("error parsing max requeue value, %w", err)
	}
	return o, nil
}
