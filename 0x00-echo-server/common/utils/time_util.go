package utils

import (
	"game/common/logger"
	"time"
)

// ProfileFunc record function time
func ProfileFunc(f func()) {
	t := NowMicro()
	f()
	delta := NowMicro() - t
	logger.Warnf("[profile] func delta:%v", delta)
}

// Now get current time
func Now() time.Time {
	t := time.Now()
	return t
}

// NowMicro get micro second
func NowMicro() int64 {
	return Now().UnixNano() / 1000
}
