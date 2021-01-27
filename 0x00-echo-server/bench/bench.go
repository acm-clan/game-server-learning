package main

import (
	"flag"
	"game/bench/client"
	"game/common/logger"
	"game/common/utils"
	"sync"
)

var clientCount = flag.Int("n", 1000, "bench client count")
var singleClientMessageCount = flag.Int64("m", 1000, "bench client message count")
var useGoroutine = flag.Bool("go", true, "bench client use goroutine")
var serverPort = flag.Int("port", 8000, "bench server port")
var logLevel = flag.String("log", "info", "bench client count")
var singleClientMessageSize = flag.Int64("s", 100, "bench client message size")

func run(isGoroutine bool) {
	logger.Infof("run go %v", isGoroutine)
	var wg sync.WaitGroup
	for i := 0; i < *clientCount; i++ {
		bc := client.NewBenchClient(int64(i+1), &wg, *singleClientMessageCount, *singleClientMessageSize)
		if isGoroutine {
			wg.Add(1)
			go bc.Start(*serverPort)
		} else {
			bc.WaitGroup = nil
			bc.Start(*serverPort)
		}
	}
	wg.Wait()
}

func main() {
	flag.Parse()
	logger.InitLogger(*logLevel)

	logger.Infof("Start benchmark:%v", *clientCount)

	utils.ProfileFunc(func() {
		run(*useGoroutine)
	})
}
