package main

import (
	"mocha/handler"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go handler.StartGrpc(&wg)
	go handler.StartRest(&wg)

	wg.Wait()
}
