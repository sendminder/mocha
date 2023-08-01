package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Conversation struct는 채팅방의 데이터 구조를 나타냅니다.
type Conversation struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

var (
	conversations []Conversation
)

// GetConversationsHandler는 모든 채팅방을 반환합니다.
func GetConversationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}

// GetConversationHandler는 특정 채팅방을 반환합니다.
func GetConversationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, conv := range conversations {
		if conv.ID == params["id"] {
			json.NewEncoder(w).Encode(conv)
			return
		}
	}
	json.NewEncoder(w).Encode(&Conversation{})
}

// CreateConversationHandler는 새로운 채팅방을 생성합니다.
func CreateConversationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var conv Conversation
	_ = json.NewDecoder(r.Body).Decode(&conv)
	conversations = append(conversations, conv)
	json.NewEncoder(w).Encode(conv)
}

func StartRest(wg *sync.WaitGroup) {
	defer wg.Done()
	router := mux.NewRouter()

	// RESTful API 핸들러 등록
	router.HandleFunc("/conversations", GetConversationsHandler).Methods("GET")
	router.HandleFunc("/conversations/{id}", GetConversationHandler).Methods("GET")
	router.HandleFunc("/conversations", CreateConversationHandler).Methods("POST")

	// 기본 HTTP 서버 시작
	fmt.Println("REST server is listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
