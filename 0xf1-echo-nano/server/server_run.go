package main

import (
	"fmt"
	"game/common/logger"
	"game/pb"
	"sync/atomic"
	"time"

	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/serialize/protobuf"
	"github.com/lonng/nano/session"
)

const (
	addr = "127.0.0.1:13250" // local address
	conc = 1000              // concurrent client count
)

//
type TestHandler struct {
	component.Base
	metrics int32
}

var totalDataSize = 0

func (h *TestHandler) AfterInit() {
	seconds := int32(2)
	ticker := time.NewTicker(time.Second * time.Duration(seconds))

	// metrics output ticker
	go func() {
		for range ticker.C {
			v := atomic.LoadInt32(&h.metrics)
			if v > 0 {
				logger.Infof("QPS %v data %v KB", v/seconds, totalDataSize/1024)
				atomic.StoreInt32(&h.metrics, 0)
				totalDataSize = 0
			}
		}
	}()
}

func (h *TestHandler) Ping(s *session.Session, data *pb.Ping) error {
	// logger.Info("recv ping")
	atomic.AddInt32(&h.metrics, 1)
	totalDataSize += len(data.Content)
	return s.Push("pong", &pb.Pong{Content: data.Content})
}

func runServer(port int) {
	comps := &component.Components{}
	comps.Register(&TestHandler{})
	nano.Listen("0.0.0.0:"+fmt.Sprintf("%v", port),
		nano.WithComponents(comps),
		nano.WithSerializer(protobuf.NewSerializer()),
	)
}
