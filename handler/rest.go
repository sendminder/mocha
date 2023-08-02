package handler

import (
	"fmt"
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

	// 기본 HTTP 서버 시작
	fmt.Println("REST server is listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
