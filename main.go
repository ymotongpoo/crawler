package main

import (
	"time"

	"crawler"
)

var MaxBoardWorkers = 2

func main() {
	for {
		crawler.Run(MaxBoardWorkers)
		time.Sleep(time.Minute * 15)
	}
}
