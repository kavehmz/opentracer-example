package main

import (
	"time"
)

func main() {
	go serveUI()
	go serveRPC()
	time.Sleep(time.Second * 60)
}
