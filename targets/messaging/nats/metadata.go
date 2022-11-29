package nats

import (
	"fmt"

	"github.com/kubemq-io/kubemq-targets/types"
)

type metadata struct {
	subject string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.subject, err = meta.MustParseString("subject")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing subject name, %w", err)
	}
	return m, nil
}
