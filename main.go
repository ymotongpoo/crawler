package main

import (
	"time"

	"crawler/crawling"
)

var MaxBoardWorkers = 2

func main() {
	for {
		crawling.Run(MaxBoardWorkers)
		time.Sleep(time.Minute * 15)
	}
}
