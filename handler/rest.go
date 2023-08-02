package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mocha/db"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Conversation struct는 채팅방의 데이터 구조를 나타냅니다.
type CreateConversation struct {
	Name        string  `json:"name,omitempty"`
	HostUserId  int64   `json:"host_user_id"`
	JoinedUsers []int64 `json:"joined_users"`
}

// GetConversationsHandler는 해당 유저의 모든 채팅방을 반환합니다.
func GetUserConversationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		// conversationIDStr이 올바른 int64로 변환되지 않은 경우 에러 처리
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
		return
	}

	conversations, err := db.GetUserConversations(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 레코드를 찾지 못한 경우 404 에러 반환
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Conversation not found"})
			return
		}
		// 다른 에러가 발생한 경우 500 에러 반환
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get conversation"})
		return
	}

	json.NewEncoder(w).Encode(conversations)
}

// GetConversationHandler는 특정 채팅방을 반환합니다.
func GetConversationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	conversationId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		// conversationIDStr이 올바른 int64로 변환되지 않은 경우 에러 처리
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid conversation ID"})
		return
	}

	conversation, err := db.GetConversationByID(conversationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 레코드를 찾지 못한 경우 404 에러 반환
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Conversation not found"})
			return
		}
		// 다른 에러가 발생한 경우 500 에러 반환
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get conversation"})
		return
	}

	json.NewEncoder(w).Encode(conversation)
}

// CreateConversationHandler는 새로운 채팅방을 생성합니다.
func CreateConversationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cc CreateConversation
	_ = json.NewDecoder(r.Body).Decode(&cc)

	var conv = db.Conversation{
		Type:          "dm",
		Name:          cc.Name,
		HostUserID:    cc.HostUserId,
		LastMessageID: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	createdConv, err := db.CreateConversation(&conv)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create conversation"})
		return
	}

	json.NewEncoder(w).Encode(createdConv)

	// conversation user 생성
	for _, value := range cc.JoinedUsers {
		var cuser = db.ConversationUser{
			ConversationID:    createdConv.ID,
			UserID:            value,
			LastSeenMessageID: 0,
		}
		err = db.CreateConversationUser(&cuser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create conversation_user"})
			return
		}
	}
}

func StartRest(wg *sync.WaitGroup) {
	defer wg.Done()
	router := mux.NewRouter()

	// RESTful API 핸들러 등록
	router.HandleFunc("/user_conversations/{id}", GetUserConversationsHandler).Methods("GET")
	router.HandleFunc("/conversations/{id}", GetConversationHandler).Methods("GET")
	router.HandleFunc("/conversations", CreateConversationHandler).Methods("POST")

	// 기본 HTTP 서버 시작
	fmt.Println("REST server is listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
