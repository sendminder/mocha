package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"mocha/internal/types"
)

type MessageRestHandler interface {
	GetMessages(w http.ResponseWriter, r *http.Request)
}

func (s *restServer) GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	channelId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		// channelIDStr이 올바른 int64로 변환되지 않은 경우 에러 처리
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid channel Id"})
		return
	}

	messages, err := s.mdb.GetMessagesByChannelID(channelId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 레코드를 찾지 못한 경우 404 에러 반환
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "messages not found"})
			return
		}
		// 다른 에러가 발생한 경우 500 에러 반환
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get messages"})
		return
	}

	json.NewEncoder(w).Encode(map[string][]types.Message{"messages": messages})
}
