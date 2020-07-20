package binding

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-targets/types"
)

var logLevelMap = map[string]string{
	"debug": "debug",
	"info":  "info",
	"error": "error",
	"":      "",
}

type options struct {
	logLevel string
}

func parseOptions(m types.Metadata) (options, error) {
	if m == nil {
		return options{}, nil
	}
	o := options{}
	var err error
	o.logLevel, err = m.ParseStringMap("log_level", logLevelMap)
	if err != nil {
		return options{}, fmt.Errorf("invalid log level property, %w", err)
	}
	return o, nil
}
