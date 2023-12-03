package main

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"sync"

	"github.com/spf13/viper"
	"mocha/internal/cache"
	"mocha/internal/config"
	"mocha/internal/db"
	"mocha/internal/handler/grpc"
	"mocha/internal/handler/rest"
)

var (
	errFailedToStartServer = errors.New("failed to start server")
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
	if rdb == nil {
		return errFailedToStartServer
	}
	mdb := db.NewDynamoDatabse(config.GetString("dynamo.host")+":"+strconv.Itoa(config.GetInt("dynamo.port")),
		config.GetString("dynamo.region"), config.GetString("dynamo.table.message"))
	if mdb == nil {
		return errFailedToStartServer
	}

	redisCache := cache.NewRedisCache(ctx, config.GetString("cache.host")+":"+strconv.Itoa(config.GetInt("cache.port")))
	if redisCache == nil {
		return errFailedToStartServer
	}

	messageServer := grpc.NewGrpcServer(config.GetInt("mocha.port"), mdb, rdb, redisCache)
	if messageServer == nil {
		return errFailedToStartServer
	}
	go messageServer.Start(&wg)

	restServer := rest.NewRestServer(rdb, mdb, redisCache, config.GetInt("rest.port"))
	if restServer == nil {
		return errFailedToStartServer
	}
	go restServer.Start(&wg)

	wg.Wait()
	return nil
}
