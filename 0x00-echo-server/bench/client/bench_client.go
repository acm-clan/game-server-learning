package client

import (
	"bufio"
	"fmt"
	"game/common/logger"
	"io"
	"net"
	"os"
	"sync"
)

// BenchClient send message to server
type BenchClient struct {
	ClientID     int64
	MessageCount int64
	WaitGroup    *sync.WaitGroup
}

func (bc *BenchClient) bench(port int) {
	conn, err := net.Dial("tcp", ":"+fmt.Sprint(port))
	if err != nil {
		logger.Error("ERROR", err)
		os.Exit(1)
	}

	response := bufio.NewReader(conn)
	for i := 0; i < int(bc.MessageCount); i++ {
		conn.Write([]byte("hello world\n"))
		serverLine, err := response.ReadBytes(byte('\n'))
		switch err {
		case nil:
			logger.Debug("[echo] " + string(serverLine))
		case io.EOF:
			os.Exit(0)
		default:
			logger.Error("ERROR", err)
			os.Exit(2)
		}
	}
}

// Start start a bench client
func (bc *BenchClient) Start(port int) {
	logger.Infof("[client] start %v", bc.ClientID)
	bc.bench(port)
	if bc.WaitGroup != nil {
		bc.WaitGroup.Done()
	}
}
