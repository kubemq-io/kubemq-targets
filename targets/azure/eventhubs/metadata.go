package eventhubs

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

const (
	DefaultProperties   = ""
	DefaultPartitionKey = ""
)

var methodsMap = map[string]string{
	"send":       "send",
	"send_batch": "send_batch",
}

type metadata struct {
	method       string
	partitionKey string
	properties   map[string]interface{}
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidMethodTypes(methodsMap)
	}
	m.properties, err = meta.MustParseInterfaceMap("properties")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing properties, %w", err)
	}
	m.partitionKey = meta.ParseString("partition_key", DefaultPartitionKey)
	return m, nil
}
