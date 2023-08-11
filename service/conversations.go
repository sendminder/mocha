package service

import (
	"encoding/json"
	"errors"
	"mocha/cache"
	"mocha/db"
	"mocha/types"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// GetConversationsHandler는 해당 유저의 모든 채팅방을 반환합니다.
func GetUserConversationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		// conversationIDStr이 올바른 int64로 변환되지 않은 경우 에러 처리
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user Id"})
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

	json.NewEncoder(w).Encode(map[string][]types.Conversation{"conversations": conversations})
}

// GetConversationHandler는 특정 채팅방을 반환합니다.
func GetConversationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	conversationId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		// conversationIDStr이 올바른 int64로 변환되지 않은 경우 에러 처리
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid conversation Id"})
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

	json.NewEncoder(w).Encode(map[string]types.Conversation{"conversation": *conversation})
}

// CreateConversationHandler는 새로운 채팅방을 생성합니다.
func CreateConversationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cc types.CreateConversation
	_ = json.NewDecoder(r.Body).Decode(&cc)

	var conv = types.Conversation{
		Type:            "dm",
		Name:            cc.Name,
		HostUserId:      cc.HostUserId,
		LastMessageId:   0,
		LastDecryptedId: 0,
		CreatedAt:       time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:       time.Now().UTC().Format(time.RFC3339),
	}
	createdConv, err := db.CreateConversation(&conv)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create conversation"})
		return
	}

	json.NewEncoder(w).Encode(map[string]types.Conversation{"conversation": *createdConv})

	// conversation user 생성
	for _, value := range cc.JoinedUsers {
		var cuser = types.ConversationUser{
			ConversationId:    createdConv.Id,
			UserId:            value,
			LastSeenMessageId: 0,
		}
		err = db.CreateConversationUser(&cuser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create conversation_user"})
			return
		}
	}

	// redis set
	cache.SetJoinedUsers(createdConv.Id, cc.JoinedUsers)
}
