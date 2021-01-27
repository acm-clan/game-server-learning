package main

import (
	"bufio"
	"game/common/logger"
	"game/common/utils"
	"math/rand"
	"net"
	"os"
	"time"
)

func handleConnection(c net.Conn) {
	logger.Infof("Serving %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			logger.Error(err)
			return
		}

		result := string(netData)

		logger.Debug(utils.DumpString(result))
		c.Write([]byte(result))
	}

	c.Close()
}

func main() {
	logger.InitLogger()

	arguments := os.Args
	if len(arguments) == 1 {
		logger.Error("Please provide a port number!")
		return
	}

	port := ":" + arguments[1]
	listener, err := net.Listen("tcp4", port)

	if err != nil {
		logger.Error(err)
		return
	}

	defer listener.Close()
	rand.Seed(time.Now().Unix())

	logger.Infof("Echo server start accept")

	for {
		c, err := listener.Accept()
		if err != nil {
			logger.Error(err)
			return
		}
		go handleConnection(c)
	}
}
