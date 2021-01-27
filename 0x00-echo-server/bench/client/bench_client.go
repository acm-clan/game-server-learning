package client

import (
	"game/common/logger"
	"sync"
)

// BenchClient send message to server
type BenchClient struct {
	ClientID     int64
	MessageCount int64
	WaitGroup    *sync.WaitGroup
}

// Start start a bench client
func (bc *BenchClient) Start() {
	logger.Infof("[client] start %v", bc.ClientID)
	if bc.WaitGroup != nil {
		bc.WaitGroup.Done()
	}
}
