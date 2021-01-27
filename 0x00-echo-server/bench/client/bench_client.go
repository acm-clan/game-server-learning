package client

import (
	"bufio"
	"fmt"
	"game/common/logger"
	"game/common/utils"
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
	connection, err := net.Dial("tcp", ":"+fmt.Sprint(port))
	if err != nil {
		logger.Error("ERROR", err)
		os.Exit(1)
	}

	response := bufio.NewReader(connection)

	for i := 0; i < int(bc.MessageCount); i++ {
		msg := utils.GenerateString(100)
		logger.Debug("[bench] msg: ", msg)
		_, err := connection.Write([]byte(msg))

		if err != nil {
			panic(err)
		}

		_, err = connection.Write([]byte("\n"))

		if err != nil {
			panic(err)
		}

		serverLine, err := response.ReadBytes(byte('\n'))
		switch err {
		case nil:
			logger.Debug("[echo] " + string(serverLine))
		default:
			logger.Error("ERROR", err)
			os.Exit(0)
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
