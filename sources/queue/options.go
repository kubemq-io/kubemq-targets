package queue

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
	"github.com/nats-io/nuid"
)

type options struct {
	host            string
	port            int
	clientId        string
	authToken       string
	channel         string
	responseChannel string
	concurrency     int
	batchSize       int
	waitTimeout     int
}

func parseOptions(cfg config.Spec) (options, error) {
	m := options{}
	var err error
	m.host = cfg.ParseString("host", defaultHost)

	m.port, err = cfg.ParseIntWithRange("port", defaultPort, 1, 65535)
	if err != nil {
		return options{}, fmt.Errorf("error parsing port value, %w", err)
	}

	m.authToken = cfg.ParseString("auth_token", "")

	m.clientId = cfg.ParseString("client_id", nuid.Next())

	m.channel, err = cfg.MustParseString("channel")
	if err != nil {
		return options{}, fmt.Errorf("error parsing channel value, %w", err)
	}
	m.responseChannel = cfg.ParseString("response_channel", "")

	m.concurrency, err = cfg.ParseIntWithRange("concurrency", 1, 1, 100)
	if err != nil {
		return options{}, fmt.Errorf("error parsing concurrency value, %w", err)
	}

	m.batchSize, err = cfg.ParseIntWithRange("batch_size", defaultBatchSize, 1, 1024)
	if err != nil {
		return options{}, fmt.Errorf("error parsing batch size value, %w", err)
	}
	m.waitTimeout, err = cfg.ParseIntWithRange("wait_timeout", defaultWaitTimeout, 1, 24*60*60)
	if err != nil {
		return options{}, fmt.Errorf("error parsing wait timeout value, %w", err)
	}

	return m, nil
}
