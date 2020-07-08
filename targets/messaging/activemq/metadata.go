package activemq

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-target-connectors/types"
)

type metadata struct {
	destination string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.destination, err = meta.MustParseString("destination")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing destination name, %w", err)
	}
	return m, nil
}
