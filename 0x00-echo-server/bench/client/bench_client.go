package client

import (
	"bufio"
	"fmt"
	"game/common/logger"
	"game/common/utils"
	"net"
	"os"
)

func (bc *BenchClient) benchSync(ip string, port int) {
	connection, err := net.Dial("tcp", ip+":"+fmt.Sprint(port))
	if err != nil {
		logger.Error("ERROR", err)
		os.Exit(1)
	}

	response := bufio.NewReader(connection)

	for i := 0; i < int(bc.MessageCount); i++ {
		msg := utils.GenerateString(int(bc.MessageSize))
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
func (bc *BenchClient) StartSync(ip string, port int) {
	logger.Debugf("[client] start %v", bc.ClientID)
	bc.benchSync(ip, port)
	if bc.WaitGroup != nil {
		bc.WaitGroup.Done()
	}
}
