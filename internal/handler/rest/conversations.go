package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"mocha/internal/types"
)

type ChannelHandler interface {
	GetUserChannels(w http.ResponseWriter, r *http.Request)
	GetChannel(w http.ResponseWriter, r *http.Request)
	CreateChannel(w http.ResponseWriter, r *http.Request)
}

// GetChannelsHandler는 해당 유저의 모든 채팅방을 반환합니다.
func (s *restServer) GetUserChannels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		// channelIDStr이 올바른 int64로 변환되지 않은 경우 에러 처리
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user Id"})
		return
	}

	channels, err := s.rdb.GetUserChannels(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 레코드를 찾지 못한 경우 404 에러 반환
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Channel not found"})
			return
		}
		// 다른 에러가 발생한 경우 500 에러 반환
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get channel"})
		return
	}

	json.NewEncoder(w).Encode(map[string][]types.Channel{"channels": channels})
}

// GetChannelHandler는 특정 채팅방을 반환합니다.
func (s *restServer) GetChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	channelId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		// channelIDStr이 올바른 int64로 변환되지 않은 경우 에러 처리
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid channel Id"})
		return
	}

	channel, err := s.rdb.GetChannelByID(channelId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 레코드를 찾지 못한 경우 404 에러 반환
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Channel not found"})
			return
		}
		// 다른 에러가 발생한 경우 500 에러 반환
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get channel"})
		return
	}

	json.NewEncoder(w).Encode(map[string]types.Channel{"channel": *channel})
}

// CreateChannelHandler는 새로운 채팅방을 생성합니다.
func (s *restServer) CreateChannel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cc types.CreateChannel
	_ = json.NewDecoder(r.Body).Decode(&cc)

	var channel = types.Channel{
		Type:            "dm",
		Name:            cc.Name,
		HostUserId:      cc.HostUserId,
		LastMessageId:   0,
		LastDecryptedId: 0,
		CreatedAt:       time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:       time.Now().UTC().Format(time.RFC3339),
	}
	createdChannel, err := s.rdb.CreateChannel(&channel)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create channel"})
		return
	}

	json.NewEncoder(w).Encode(map[string]types.Channel{"channel": *createdChannel})

	// channel user 생성
	for _, value := range cc.JoinedUsers {
		var cuser = types.ChannelUser{
			ChannelId:         createdChannel.Id,
			UserId:            value,
			LastSeenMessageId: 0,
		}
		err = s.rdb.CreateChannelUser(&cuser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create channel_user"})
			return
		}
	}

	// redis set
	s.cache.SetJoinedUsers(createdChannel.Id, cc.JoinedUsers)
}
