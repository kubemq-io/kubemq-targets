// forked from https://github.com/uber-go/ratelimit
package clock

import "time"

type Clock interface {
	AfterFunc(d time.Duration, f func())
	Now() time.Time
	Sleep(d time.Duration)
}
