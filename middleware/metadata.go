package middleware

import "github.com/kubemq-io/kubemq-targets/types"

type MetadataMiddleware struct {
	Metadata map[string]string
}

func NewMetadataMiddleware(meta types.Metadata) (*MetadataMiddleware, error) {
	mm := &MetadataMiddleware{}
	var err error
	mm.Metadata, err = meta.MustParseJsonMap("metadata")
	if err != nil {
		mm.Metadata = map[string]string{}
	}
	return mm, nil
}
