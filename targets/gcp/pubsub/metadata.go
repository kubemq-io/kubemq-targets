package pubsub

import (
	"fmt"
	"github.com/kubemq-io/kubemq-targets/types"
)

type metadata struct {
	topicID string
	tags    map[string]string
}

func parseMetadata(meta types.Metadata, opts options) (metadata, error) {
	m := metadata{}
	m.tags = make(map[string]string)
	var err error
	m.topicID, err = meta.MustParseString("topic_id")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing topic_id, %w", err)
	}
	tags, err := meta.MustParseJsonMap("tags")
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing tags, %w", err)
	}
	m.tags = tags
	return m, nil
}
