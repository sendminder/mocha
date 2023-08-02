package main

import (
	"mocha/db"
	"mocha/handler"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	// go db.ConnectPostgresql(&wg)
	go db.ConnectGorm(&wg)
	go handler.StartGrpc(&wg)
	go handler.StartRest(&wg)

	wg.Wait()
}
