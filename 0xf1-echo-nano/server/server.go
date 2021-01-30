package main

import (
	"bufio"
	"flag"
	"game/common/logger"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
)

var logLevel = flag.String("log", "info", "log level")
var serverPort = flag.Int("port", 8000, "server port")
var useSync = flag.Bool("sync", true, "server use sync")
var echoBack = flag.Bool("echo", true, "server echo")

func handleConnectionSync(c net.Conn) {
	logger.Debugf("Serving %s\n", c.RemoteAddr().String())

	msgCount := 0
	r := bufio.NewReader(c)

	for {
		msg, err := r.ReadString('\n')
		if err == io.EOF {
			logger.Debugf("client closed: %v", err)
			break
		}

		if err != nil {
			logger.Errorf("read error: %v", err)
			break
		}

		msgCount++

		logger.Debugf("[server] %v recv: %v", msgCount, msg)

		if *echoBack {
			_, err = c.Write([]byte(msg))

			if err != nil {
				logger.Errorf("write error: %v", err)
				break
			}
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

	runServer(*serverPort)
}
