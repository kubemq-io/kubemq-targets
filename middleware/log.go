package middleware

import (
	"github.com/kubemq-hub/kubemq-targets/pkg/logger"
)

type LogLevelType string

const (
	LogLevelTypeNoLog LogLevelType = ""
	LogLevelTypeDebug LogLevelType = "debug"
	LogLevelTypeInfo  LogLevelType = "info"
	LogLevelTypeError LogLevelType = "error"
)

type LogMiddleware struct {
	minLevel LogLevelType
	*logger.Logger
}

func NewLogMiddleware(name, level string) *LogMiddleware {
	lm := &LogMiddleware{
		minLevel: LogLevelTypeNoLog,
		Logger:   logger.NewLogger(name),
	}
	switch level {
	case "debug":
		lm.minLevel = LogLevelTypeDebug
	case "info":
		lm.minLevel = LogLevelTypeInfo
	case "error":
		lm.minLevel = LogLevelTypeError
	}
	return lm
}
