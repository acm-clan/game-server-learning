package client

import (
	"bufio"
	"fmt"
	"game/common/logger"
	"game/common/utils"
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

	for i := 0; i < int(bc.MessageCount); i++ {
		msg := utils.GenerateString(int(bc.MessageSize)) + "\n"
		logger.Debug("[bench] msg: ", msg)

		_, err := bc.Connection.Write([]byte(msg))
		if err != nil {
			panic(err)
		}

		ret, err := response.ReadBytes(byte('\n'))

		logger.Debugf("[echo] %v %v ", i+1, string(ret))

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
