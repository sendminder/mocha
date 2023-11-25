package main

import (
	"context"
	"log/slog"
	"strconv"
	"sync"

	"github.com/spf13/viper"
	cache "mocha/internal/cache"
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
	mdb := db.NewDynamoDatabse(config.GetString("dynamo.host")+":"+strconv.Itoa(config.GetInt("dynamo.port")),
		config.GetString("dynamo.region"), config.GetString("dynamo.table.message"))
	cache := cache.NewRedisCache(config.GetString("cache.host") + ":" + strconv.Itoa(config.GetInt("cache.port")))

	messageServer := grpc.NewGrpcServer(config.GetInt("mocha.port"), mdb, rdb, cache)
	go messageServer.Start(&wg)

	restServer := rest.NewRestServer(rdb, mdb, cache)
	go restServer.Start(&wg)

	wg.Wait()

	return nil
}
