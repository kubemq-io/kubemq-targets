package sqs

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/config"
)

const (
	DefaultRetries    = 0
	DefaultDelay      = 10
	DefaultMaxReceive = 0
	DefaultToken      = ""
	DefaultDeadLetter = ""
)

type options struct {
	sqsKey          string
	sqsSecretKey    string
	retries         int
	region          string
	maxReceiveCount int
	deadLetterQueue string
	token           string
	defaultDelay    int
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.sqsKey, err = cfg.Properties.MustParseString("aws_key")
	if err != nil {
		return options{}, fmt.Errorf("error sqsKey , %w", err)
	}

	o.sqsSecretKey, err = cfg.Properties.MustParseString("aws_secret_key")
	if err != nil {
		return options{}, fmt.Errorf("error sqsSecretKey , %w", err)
	}

	o.retries = cfg.Properties.ParseInt("retries", DefaultRetries)

	o.region, err = cfg.Properties.MustParseString("region")
	if err != nil {
		return options{}, fmt.Errorf("error parsing region , %w", err)
	}

	o.defaultDelay = DefaultDelay
	o.maxReceiveCount = cfg.Properties.ParseInt("max_receive", DefaultMaxReceive)
	o.deadLetterQueue = cfg.Properties.ParseString("dead_letter", DefaultDeadLetter)
	o.token = cfg.Properties.ParseString("token", DefaultToken)

	return o, nil
}
