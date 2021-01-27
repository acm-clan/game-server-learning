package utils

import (
	"game/common/logger"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// ProfileFunc record function time
func ProfileFunc(f func()) {
	t := NowMicro()
	f()
	delta := NowMicro() - t

	p := message.NewPrinter(language.English)
	s := p.Sprintf("%v", delta)

	logger.Warnf("[profile] func delta: %v micro second", s)
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