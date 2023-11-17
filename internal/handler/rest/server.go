package rest

import (
	"log/slog"
	"mocha/internal/cache"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"mocha/internal/db"
)

type RestServer interface {
	ConversationHandler
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
}

func (s *restServer) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	router := mux.NewRouter()

	// RESTful API 핸들러 등록
	router.HandleFunc("/user_conversations/{id}", s.GetUserConversations).Methods("GET")
	router.HandleFunc("/conversations/{id}", s.GetConversation).Methods("GET")
	router.HandleFunc("/conversations", s.CreateConversation).Methods("POST")

	router.HandleFunc("/conversations/{id}/messages", s.GetMessages).Methods("GET")

	router.HandleFunc("/users/{id}", s.GetUser).Methods("GET")
	router.HandleFunc("/users", s.CreateUser).Methods("POST")
	router.HandleFunc("/users/login", s.LoginUser).Methods("POST")

	router.HandleFunc("/devices/{id}", s.GetDevice).Methods("GET")
	router.HandleFunc("/devices", s.CreateDevice).Methods("POST")

	slog.Info("REST server is listening on port 8000...")
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		return
	}
}

func NewRestServer(rdb db.RelationalDatabase, mdb db.DynamoDatabase, cache cache.RedisCache) RestServer {
	return &restServer{
		rdb:   rdb,
		mdb:   mdb,
		cache: cache,
	}
}
