package mqtt

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

type metadata struct {
	topic string
	qos   int
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.topic, err = meta.MustParseString("topic")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing topic name, %w", err)
	}
	m.qos, err = meta.ParseIntWithRange("qos", 0, 0, 2)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing qos, %w", err)
	}
	return m, nil
}
