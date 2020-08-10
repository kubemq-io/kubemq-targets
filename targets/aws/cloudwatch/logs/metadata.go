package logs

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

const (
	defaultLimit = 100
)

type metadata struct {
	method string

	limit          int64
	logGroupName   string
	logStreamName  string
	logGroupPrefix string
}

var methodsMap = map[string]string{
	"create_log_event_stream":   "create_log_event_stream",
	"describe_log_event_stream": "describe_log_event_stream",
	"delete_log_event_stream":   "delete_log_event_stream",
	"put_log_event":             "put_log_event",
	"get_log_event":             "get_log_event",
	"create_log_group":          "create_log_group",
	"delete_log_group":          "delete_log_group",
	"describe_log_group":        "describe_log_group",
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
	m.limit = int64(meta.ParseInt("limit", defaultLimit))
	if err != nil {
		return metadata{}, fmt.Errorf(getValidMethodTypes())
	}
	if m.method == "describe_log_group" {
		m.logGroupPrefix, err = meta.MustParseString("log_group_prefix")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing log_group_prefix, %w", err)
		}
	} else {
		m.logGroupName, err = meta.MustParseString("log_group_name")
		if err != nil {
			return metadata{}, fmt.Errorf("error parsing log_group_name, %w", err)
		}
		if m.method == "put_log_event" || m.method == "get_log_event" || m.method == "create_log_event_stream" || m.method == "delete_log_event_stream" {
			m.logStreamName, err = meta.MustParseString("log_stream_name")
			if err != nil {
				return metadata{}, fmt.Errorf("error parsing log_stream_name, %w", err)
			}
		}
	}
	return m, nil
}
