package cloudfunctions

import (
	"fmt"

	"github.com/kubemq-io/kubemq-targets/types"
)

type metadata struct {
	name     string
	project  string
	location string
}

func parseMetadata(meta types.Metadata, opts options) (metadata, error) {
	m := metadata{}

	var err error
	m.name, err = meta.MustParseString("name")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing name, %w", err)
	}
	m.project = meta.ParseString("project", "")
	m.location = meta.ParseString("location", "")

	return m, nil
}
