package main

import (
	"bufio"
	"flag"
	"fmt"
	"game/common/logger"
	"game/server/client"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var logLevel = flag.String("log", "info", "log level")
var serverPort = flag.Int("port", 8000, "server port")
var useSync = flag.Bool("sync", true, "server use sync")

func handleConnectionSync(c net.Conn) {
	logger.Debugf("Serving %s\n", c.RemoteAddr().String())

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			logger.Debug(err)
			break
		}

		result := string(netData)

		_, err = c.Write([]byte(result))

		if err != nil {
			panic(err)
		}
	}

	c.Close()
}

func serviceProfile() {
	logger.Info("Start profile 8001")
	http.ListenAndServe("0.0.0.0:8001", nil)

}

func startProfile() {
	go serviceProfile()
}

func main() {
	flag.Parse()
	logger.InitLogger(*logLevel)
	startProfile()

	port := ":" + fmt.Sprint(*serverPort)
	listener, err := net.Listen("tcp4", port)

	if err != nil {
		logger.Error(err)
		return
	}

	defer listener.Close()
	rand.Seed(time.Now().Unix())

	logger.Infof("Echo server start accept sync %v", *useSync)

	for {
		c, err := listener.Accept()
		if err != nil {
			logger.Error(err)
			return
		}

		if *useSync {
			go handleConnectionSync(c)
		} else {
			cl := client.NewClient(c)
			go cl.HandleConnection()
		}

	}
}
