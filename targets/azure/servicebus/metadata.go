package servicebus

import (
	"time"

	"github.com/Azure/azure-service-bus-go"
	"github.com/kubemq-io/kubemq-targets/types"
)

const (
	DefaultTimeToLive   = 1000000000
	DefaultContentType  = ""
	DefaultLabel        = ""
	DefaultMaxBatchSize = 1024
)

var methodsMap = map[string]string{
	"send":       "send",
	"send_batch": "send_batch",
}

type metadata struct {
	method       string
	label        string
	contentType  string
	maxBatchSize servicebus.MaxMessageSizeInBytes
	timeToLive   time.Duration
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		m.method = "send"
	}
	m.timeToLive = meta.ParseTimeDuration("time_to_live", DefaultTimeToLive)
	m.contentType = meta.ParseString("content_type", DefaultContentType)
	m.label = meta.ParseString("label", DefaultLabel)
	maxBatchSize := meta.ParseInt("max_batch_size", DefaultMaxBatchSize)
	m.maxBatchSize = servicebus.MaxMessageSizeInBytes(maxBatchSize)
	return m, nil
}
