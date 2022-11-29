package queue

import (
	"fmt"
	"time"

	"github.com/kubemq-io/kubemq-targets/types"
)

const (
	DefaultMaxMessages       = 32
	DefaultVisibilityTimeout = 100000000
	DefaultTimeToLive        = 100000000
)

var methodsMap = map[string]string{
	"create":             "create",
	"get_messages_count": "get_messages_count",
	"peek":               "peek",
	"push":               "push",
	"pop":                "pop",
	"delete":             "delete",
}

type metadata struct {
	method            string
	queueName         string
	serviceUrl        string
	maxMessages       int32
	visibilityTimeout time.Duration
	timeToLive        time.Duration
	queueMetadata     map[string]string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}
	m.queueName, err = meta.MustParseString("queue_name")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing queue_name , %w", err)
	}
	m.serviceUrl, err = meta.MustParseString("service_url")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing service_url , %w", err)
	}

	queueMetadata, err := meta.MustParseJsonMap("queue_metadata")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing queue_metadata, %w", err)
	} else {
		m.queueMetadata = queueMetadata
	}

	m.maxMessages = int32(meta.ParseInt("max_messages", DefaultMaxMessages))
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing max_messages, %w", err)
	}

	m.visibilityTimeout = meta.ParseTimeDuration("visibility_timeout", DefaultVisibilityTimeout)
	m.timeToLive = meta.ParseTimeDuration("time_to_live", DefaultTimeToLive)
	m.timeToLive = m.timeToLive * time.Second
	return m, nil
}
