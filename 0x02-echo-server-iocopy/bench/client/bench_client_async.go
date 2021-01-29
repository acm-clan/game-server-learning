package client

import (
	"bufio"
	"fmt"
	"game/common/logger"
	"game/common/utils"
	"net"
	"sync"
)

type Message struct {
	data string
}

// BenchClient send message to server
type BenchClient struct {
	ClientID          int64
	MessageCount      int64
	MessageWriteCount int64
	MessageReadCount  int64
	WaitGroup         *sync.WaitGroup
	MessageSize       int64
	chClose           chan struct{}
	Connection        net.Conn
	Sync              bool
}

func NewBenchClient(CID int64, wg *sync.WaitGroup, messageCount int64, messageSize int64) *BenchClient {
	return &BenchClient{
		ClientID:          CID,
		MessageCount:      messageCount,
		MessageReadCount:  0,
		MessageWriteCount: 0,
		MessageSize:       messageSize,
		WaitGroup:         wg,
		chClose:           make(chan struct{}),
		Sync:              true,
	}
}

func (bc *BenchClient) BeginWrite() {
	for i := 0; i < int(bc.MessageCount); i++ {
		msg := utils.GenerateString(int(bc.MessageSize)) + "\n"
		size, err := bc.Connection.Write([]byte(msg))

		logger.Debugf("[bench] %v write:%v %v", i+1, utils.DumpString(msg), size)

		if err != nil {
			logger.Infof("write error: %v", err)
			break
		}

		if size != int(bc.MessageSize+1) {
			logger.Infof("write size error %v", size)
			break
		}
	}

	logger.Debug("bench client write exit normally")
}

func (bc *BenchClient) benchAsync(ip string, port int) {
	connection, err := net.Dial("tcp", ip+":"+fmt.Sprint(port))
	if err != nil {
		logger.Error("[bench] error connect server: ", err)
		return
	}

	bc.Connection = connection

	go bc.BeginWrite()

	response := bufio.NewReader(connection)

	for {
		msg, err := response.ReadBytes(byte('\n'))

		bc.MessageReadCount++

		logger.Debugf("[echo] %v %v ", bc.MessageReadCount, string(msg))

		if err != nil {
			logger.Errorf("read error: %v", err)
			close(bc.chClose)
			break
		}

		if bc.MessageReadCount == bc.MessageCount {
			break
		}
	}

	logger.Debugf("[bench] read finished")
	bc.Connection.Close()
}

// Start start a bench client
func (bc *BenchClient) StartAsync(ip string, port int) {
	logger.Debugf("[client] start %v", bc.ClientID)
	bc.benchAsync(ip, port)
	if bc.WaitGroup != nil {
		bc.WaitGroup.Done()
	}
}

// Start start a bench client
func (bc *BenchClient) Start(ip string, port int) {
	if bc.Sync {
		bc.StartSync(ip, port)
	} else {
		bc.StartAsync(ip, port)
	}
}
