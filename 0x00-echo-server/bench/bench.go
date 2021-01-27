package main

import (
	"flag"
	"game/bench/bench"
	"game/common/logger"
)

var clientCount = flag.Int("n", 1000, "bench client count")

func main() {
	logger.InitLogger()
	flag.Parse()

	logger.Infof("Start benchmark:%v", *clientCount)

	for i := 0; i < *clientCount; i++ {
		bc := &bench.BenchClient{
			ClientID: int64(i + 1),
		}
		bc.Start()
	}
}
