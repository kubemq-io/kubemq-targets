package eventhubs

import (
	"github.com/kubemq-hub/kubemq-targets/types"
)

const (
	DeleteProperties   = ""
	DeletePartitionKey = ""
)

var methodsMap = map[string]string{
	"send":       "send",
	"send_batch": "send_batch",
}

type metadata struct {
	method       string
	partitionKey string
	properties   string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}
	m.properties = meta.ParseString("properties", DeleteProperties)
	m.partitionKey = meta.ParseString("partition_key", DeletePartitionKey)
	return m, nil
}
