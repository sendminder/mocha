package handler

import (
	"log"
	"mocha/service"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func StartRest(wg *sync.WaitGroup) {
	defer wg.Done()
	router := mux.NewRouter()

	// RESTful API 핸들러 등록
	router.HandleFunc("/user_conversations/{id}", service.GetUserConversationsHandler).Methods("GET")
	router.HandleFunc("/conversations/{id}", service.GetConversationHandler).Methods("GET")
	router.HandleFunc("/conversations", service.CreateConversationHandler).Methods("POST")

	router.HandleFunc("/conversations/{id}/messages", service.GetMessages).Methods("GET")

	router.HandleFunc("/users/{id}", service.GetUserHandler).Methods("GET")
	router.HandleFunc("/users", service.CreateUserHandler).Methods("POST")

	router.HandleFunc("/devices/{id}", service.GetDeviceHandler).Methods("GET")
	router.HandleFunc("/devices", service.CreateDeviceHandler).Methods("POST")

	log.Println("REST server is listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
