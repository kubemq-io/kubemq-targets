package events

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

const (
	defaultDetail     = ""
	defaultDetailType = ""
	defaultSource     = ""
)

type metadata struct {
	method string
	rule   string

	detail     string
	detailType string
	source     string
}

var methodsMap = map[string]string{
	"put_targets": "put_targets",
	"send_event":  "send_event",
}

func getValidMethodTypes() string {
	s := "invalid method type, method type should be one of the following:"
	for k := range methodsMap {
		s = fmt.Sprintf("%s :%s,", s, k)
	}
	return s
}

func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, fmt.Errorf(getValidMethodTypes())
	}
	m.detail = meta.ParseString("detail", defaultDetail)
	m.detailType = meta.ParseString("detail_type", defaultDetailType)
	m.source = meta.ParseString("source", defaultSource)
	if m.method == "put_targets" {
		m.rule, err = meta.MustParseString("rule")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing rule, %w", err)
		}
	}
	return m, nil
}
