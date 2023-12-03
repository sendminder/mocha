package rest

import (
	"log/slog"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
	"mocha/internal/cache"
	"mocha/internal/db"
)

type RestServer interface {
	ChannelHandler
	UserHandler
	MessageRestHandler
	DeviceHandler
	Start(wg *sync.WaitGroup)
}

var _ RestServer = (*restServer)(nil)

type restServer struct {
	rdb   db.RelationalDatabase
	mdb   db.DynamoDatabase
	cache cache.RedisCache
	port  int
}

func (s *restServer) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	app := fiber.New()

	app.Get("/user_channels/:id", s.GetUserChannels())
	app.Get("/channels/:id", s.GetChannel())
	app.Post("/channels", s.CreateChannel())

	app.Get("/channels/:id/messages", s.GetMessages())

	app.Get("/users/:id", s.GetUser())
	app.Post("/users", s.CreateUser())
	app.Post("/users/login", s.LoginUser())

	app.Get("/devices/:id", s.GetDevice())
	app.Post("/devices", s.CreateDevice())

	slog.Info("REST server is listening on port", "port", strconv.Itoa(s.port))
	err := app.Listen(":" + strconv.Itoa(s.port))
	if err != nil {
		return
	}
}

func NewRestServer(rdb db.RelationalDatabase, mdb db.DynamoDatabase, cache cache.RedisCache, port int) RestServer {
	return &restServer{
		rdb:   rdb,
		mdb:   mdb,
		cache: cache,
		port:  port,
	}
}
