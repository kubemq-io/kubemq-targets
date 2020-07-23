package cloudfunctions

import (
	"fmt"

	"github.com/kubemq-hub/kubemq-targets/types"
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
		return metadata{}, fmt.Errorf("error parsing topic_id, %w", err)
	}
	m.project = meta.ParseString("project", "")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing project, %w", err)
	}
	m.location = meta.ParseString("location", "")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing location, %w", err)
	}

	return m, nil
}
