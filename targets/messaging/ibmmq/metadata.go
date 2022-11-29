//go:build container
// +build container

package ibmmq

import (
	"github.com/kubemq-io/kubemq-targets/types"
)

type metadata struct {
	dynamicQueue string
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	m.dynamicQueue = meta.ParseString("dynamic_queue", "")
	return m, nil
}
