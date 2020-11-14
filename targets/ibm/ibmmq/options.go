package ibmmq

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
)

const (
	defaultCertificateLabel = ""
	defaultTimeToLive       = 3600000
	deliveryMode            = 1
)

type options struct {
	qMName           string
	hostname         string
	portNumber       int
	channelName      string
	userName         string
	keyRepository    string
	certificateLabel string
	queueName        string
	Password         string
	deliveryMode     int
	timeToLive       int
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.qMName, err = cfg.Properties.MustParseString("queue_manager_name")
	if err != nil {
		return options{}, fmt.Errorf("error parsing qm_name, %w", err)
	}
	o.hostname, err = cfg.Properties.MustParseString("host_name")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host_name, %w", err)
	}
	o.portNumber, err = cfg.Properties.MustParseInt("port_number")
	if err != nil {
		return options{}, fmt.Errorf("error parsing port_number, %w", err)
	}
	o.channelName, err = cfg.Properties.MustParseString("channel_name")
	if err != nil {
		return options{}, fmt.Errorf("error parsing channel_name, %w", err)
	}
	o.userName, err = cfg.Properties.MustParseString("user_name")
	if err != nil {
		return options{}, fmt.Errorf("error parsing user_name, %w", err)
	}
	o.keyRepository, err = cfg.Properties.MustParseString("key_repository")
	if err != nil {
		return options{}, fmt.Errorf("error parsing key_repository, %w", err)
	}
	o.Password, err = cfg.Properties.MustParseString("password")
	if err != nil {
		return options{}, fmt.Errorf("error parsing password, %w", err)
	}
	o.queueName, err = cfg.Properties.MustParseString("queue_name")
	if err != nil {
		return options{}, fmt.Errorf("error parsing queue_name, %w", err)
	}
	o.certificateLabel = cfg.Properties.ParseString("certificate_label", defaultCertificateLabel)
	o.timeToLive = cfg.Properties.ParseInt("ttl", defaultTimeToLive)
	o.deliveryMode = cfg.Properties.ParseInt("delivery_mode", deliveryMode)

	return o, nil
}
