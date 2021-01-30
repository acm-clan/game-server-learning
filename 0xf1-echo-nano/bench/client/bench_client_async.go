package client

import (
	"bufio"
	"fmt"
	"game/common/logger"
	"game/common/utils"
	"game/pb"
	"io"
	"net"
	"os"
	"sync"
	"time"
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

func (bc *BenchClient) Write0() {
	fileName := fmt.Sprintf("files/bench.txt")

	var f *os.File
	var err error

	f, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error("file open fail", err)
		return
	}

	t := 0
	for i := 0; i < int(bc.MessageCount); i++ {
		msg := utils.GenerateString(int(bc.MessageSize)) + "\n"
		n, err := io.WriteString(f, msg)
		if err != nil {
			logger.Error("write error", err)
			return
		}
		t += n
	}

	logger.Infof("write file %v size %v KB", fileName, t/1024)
}

func (bc *BenchClient) GetReader() io.Reader {
	fileName := fmt.Sprintf("files/bench.txt")

	var f *os.File
	var err error

	f, err = os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		logger.Error("read file open fail", err)
		return nil
	}
	r := bufio.NewReader(f)
	return r
}

func (bc *BenchClient) BeginWrite() {
	//bc.Write0()
	r := bc.GetReader()
	w := bufio.NewWriter(bc.Connection)
	buf := make([]byte, 4096)
	_, err := io.CopyBuffer(w, r, buf)

	if err != nil {
		logger.Error("write error: ", err)
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
func (bc *BenchClient) Start2(ip string, port int) {
	if bc.Sync {
		bc.StartSync(ip, port)
	} else {
		bc.StartAsync(ip, port)
	}
}

// Start start a bench client
func (bc *BenchClient) Start(ip string, port int) {
	c := NewConnector()

	chReady := make(chan struct{})
	c.OnConnected(func() {
		chReady <- struct{}{}
		logger.Info("connect server")
	})

	if err := c.Start("127.0.0.1:8000"); err != nil {
		panic(err)
	}

	c.On("pong", func(data interface{}) {})

	<-chReady
	for c.chSend != nil {
		logger.Info("send notify")
		c.Notify("TestHandler.Ping", &pb.Ping{})
		time.Sleep(100 * time.Millisecond)
	}

	if bc.WaitGroup != nil {
		bc.WaitGroup.Done()
	}
}
