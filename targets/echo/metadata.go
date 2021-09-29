package echo

import (
	"github.com/kubemq-io/kubemq-targets/types"
)

type metadata struct {
	isError bool
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	m.isError = meta.ParseBool("is-error", false)
	return m, nil
}
