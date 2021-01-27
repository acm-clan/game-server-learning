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
	chWrite           chan *Message
	chClose           chan struct{}
	Connection        net.Conn
}

func NewBenchClient(CID int64, wg *sync.WaitGroup, messageCount int64, messageSize int64) *BenchClient {
	return &BenchClient{
		ClientID:          CID,
		MessageCount:      messageCount,
		MessageReadCount:  0,
		MessageWriteCount: 0,
		MessageSize:       messageSize,
		WaitGroup:         wg,
		chWrite:           make(chan *Message),
		chClose:           make(chan struct{}),
	}
}

func (bc *BenchClient) WriteRandom() {
	if bc.MessageWriteCount == bc.MessageCount {
		return
	}

	msg := utils.GenerateString(int(bc.MessageSize)) + "\n"
	logger.Debug("[bench] msg: ", msg)

	m := &Message{
		data: msg,
	}

	bc.chWrite <- m

	bc.MessageWriteCount++
}

func (bc *BenchClient) BeginWrite() {
	defer func() {
		close(bc.chWrite)
		bc.Connection.Close()
	}()

	for bc.chClose != nil {
		select {
		case data := <-bc.chWrite:
			_, err := bc.Connection.Write([]byte(data.data))
			if err != nil {
				panic(err)
			}
			break
		case <-bc.chClose:
			bc.chClose = nil
			break
		}
	}
	logger.Debug("bench client write exit normally")
}

func (bc *BenchClient) bench(ip string, port int) {
	connection, err := net.Dial("tcp", ip+":"+fmt.Sprint(port))
	if err != nil {
		logger.Error("[bench] error connect server: ", err)
		return
	}

	bc.Connection = connection

	go bc.BeginWrite()

	bc.WriteRandom()

	response := bufio.NewReader(connection)

	for {
		msg, err := response.ReadBytes(byte('\n'))

		bc.MessageReadCount++

		logger.Debugf("[echo] %v %v ", bc.MessageReadCount, string(msg))

		if err != nil || bc.MessageReadCount == bc.MessageCount {
			close(bc.chClose)
			break
		}

		bc.WriteRandom()
	}

}

// Start start a bench client
func (bc *BenchClient) Start(ip string, port int) {
	logger.Debugf("[client] start %v", bc.ClientID)
	bc.bench(ip, port)
	if bc.WaitGroup != nil {
		bc.WaitGroup.Done()
	}
}
