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

func run(isGoroutine bool) {
	logger.Infof("run go %v", isGoroutine)
	var wg sync.WaitGroup
	for i := 0; i < *clientCount; i++ {
		bc := &client.BenchClient{
			ClientID:     int64(i + 1),
			MessageCount: *singleClientMessageCount,
			WaitGroup:    &wg,
		}
		if isGoroutine {
			wg.Add(1)
			go bc.Start()
		} else {
			bc.WaitGroup = nil
			bc.Start()
		}
	}
	wg.Wait()
}

func main() {
	logger.InitLogger()
	flag.Parse()

	logger.Infof("Start benchmark:%v", *clientCount)

	utils.ProfileFunc(func() {
		run(*useGoroutine)
	})
}
