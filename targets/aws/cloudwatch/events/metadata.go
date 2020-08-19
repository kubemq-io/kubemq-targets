package events

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

const (
	defaultDetail     = ""
	defaultDetailType = ""
	defaultSource     = ""
	defaultLimit      = 10
)

type metadata struct {
	method     string
	rule       string
	limit      int64
	detail     string
	detailType string
	source     string
}

var methodsMap = map[string]string{
	"put_targets": "put_targets",
	"send_event":  "send_event",
	"list_buses":  "list_buses",
}


func parseMetadata(meta types.Metadata) (metadata, error) {
	m := metadata{}
	var err error
	m.method, err = meta.ParseStringMap("method", methodsMap)
	if err != nil {
		return metadata{}, meta.GetValidHttpMethodTypes(methodsMap)
	}
	m.limit = int64(meta.ParseInt("limit", defaultLimit))
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
