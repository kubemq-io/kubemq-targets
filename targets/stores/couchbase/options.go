package couchbase

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/config"
	"math"
)

const (
	defaultNumToReplicate = 1
	defaultNumToPersist   = 1
)

type options struct {
	url            string
	username       string
	password       string
	bucket         string
	numToReplicate int
	numToPersist   int
	collection     string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.url, err = cfg.Properties.MustParseString("url")
	if err != nil {
		return options{}, fmt.Errorf("error parsing url, %w", err)
	}
	o.username = cfg.Properties.ParseString("username", "")
	o.password = cfg.Properties.ParseString("password", "")
	o.bucket, err = cfg.Properties.MustParseString("bucket")
	if err != nil {
		return options{}, fmt.Errorf("error parsing cluster name, %w", err)
	}
	o.numToReplicate, err = cfg.Properties.ParseIntWithRange("num_to_replicate", defaultNumToReplicate, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing num to replicate, %w", err)
	}
	o.numToPersist, err = cfg.Properties.ParseIntWithRange("num_to_persist", defaultNumToPersist, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing num to persist, %w", err)
	}
	o.collection = cfg.Properties.ParseString("collection", "")
	return o, nil
}
