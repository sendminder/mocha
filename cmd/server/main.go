package main

import (
	"context"
	"log/slog"
	"sync"

	"github.com/spf13/viper"
	"mocha/internal/config"
	"mocha/internal/db"
	"mocha/internal/handler/grpc"
	"mocha/internal/handler/rest"
)

func main() {
	viper.AutomaticEnv()
	ctx := context.Background()
	if err := run(ctx); err != nil {
		test := 1
		slog.Error("failed to run", "test", test)
	}
	slog.Info("server terminated")
}

func run(ctx context.Context) error {
	config.Init()

	var wg sync.WaitGroup
	wg.Add(2)

	rdb := db.NewRelationalDatabase(&wg, 10)
	mdb := db.NewDynamoDatabse("http://localhost:8001", "us-west-2", "messages")
	messageServer := grpc.NewGrpcServer(3000, mdb, rdb)
	go messageServer.Start(&wg)

	restServer := rest.NewRestServer(rdb, mdb)
	go restServer.Start(&wg)

	wg.Wait()

	return nil
}
