package openfaas

import (
	"fmt"

	"github.com/kubemq-io/kubemq-targets/types"
)

type metadata struct {
	topic string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.topic, err = meta.MustParseString("topic")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing topic func, %w", err)
	}
	return m, nil
}
