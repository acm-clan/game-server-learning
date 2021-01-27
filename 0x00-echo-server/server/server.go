package main

import (
	"flag"
	"fmt"
	"game/common/logger"
	"game/server/client"
	"math/rand"
	"net"
	"time"
)

var logLevel = flag.String("log", "info", "log level")
var serverPort = flag.Int("port", 8000, "server port")


func handleConnectionSync(c net.Conn) {
	logger.Debugf("Serving %s\n", c.RemoteAddr().String())

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			logger.Debug(err)
			break
		}

		result := string(netData)

		logger.Debug(utils.DumpString(result))
		c.Write([]byte(result))
	}

	c.Close()
}

func main() {
	flag.Parse()
	logger.InitLogger(*logLevel)

	port := ":" + fmt.Sprint(*serverPort)
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

		cl := client.NewClient(c)
		go cl.HandleConnection()
	}
}
