package client

import (
	"bufio"
	"fmt"
	"game/common/logger"
	"net"
)

func (bc *BenchClient) benchSync(ip string, port int) {
	connection, err := net.Dial("tcp", ip+":"+fmt.Sprint(port))
	if err != nil {
		logger.Error("[bench] error connect server: ", err)
		return
	}

	bc.Connection = connection

	response := bufio.NewReader(connection)

	for {
		msg, err := response.ReadBytes(byte('\n'))

		bc.MessageReadCount++

		logger.Debugf("[echo] %v %v ", bc.MessageReadCount, string(msg))

		if err != nil || bc.MessageReadCount == bc.MessageCount {
			close(bc.chClose)
			break
		}

		_, err = bc.Connection.Write(msg)
		if err != nil {
			panic(err)
		}
	}

}

// Start start a bench client
func (bc *BenchClient) StartSync(ip string, port int) {
	logger.Debugf("[client] start %v", bc.ClientID)
	bc.benchSync(ip, port)
	if bc.WaitGroup != nil {
		bc.WaitGroup.Done()
	}
}
