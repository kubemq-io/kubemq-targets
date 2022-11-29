//go:build !container
// +build !container

package global

const (
	DefaultApiPort = 8081
	EnableLogFile  = true
	LoggerType     = "console"
)
