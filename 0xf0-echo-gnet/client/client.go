package main

import (
	"bufio"
	"game/common/logger"
	"io"
	"net"
	"os"
)

func main() {
	logger.InitLogger()

	arguments := os.Args
	if len(arguments) == 1 {
		logger.Error("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]

	conn, err := net.Dial("tcp", PORT)
	if err != nil {
		logger.Error("ERROR", err)
		os.Exit(1)
	}

	userInput := bufio.NewReader(os.Stdin)
	response := bufio.NewReader(conn)
	for {
		userLine, err := userInput.ReadBytes(byte('\n'))
		switch err {
		case nil:
			conn.Write(userLine)
			logger.Info("[send] " + string(userLine))
		case io.EOF:
			os.Exit(0)
		default:
			logger.Error("ERROR", err)
			os.Exit(1)
		}

		serverLine, err := response.ReadBytes(byte('\n'))
		switch err {
		case nil:
			logger.Info("[echo] " + string(serverLine))
		case io.EOF:
			os.Exit(0)
		default:
			logger.Error("ERROR", err)
			os.Exit(2)
		}
	}
}
