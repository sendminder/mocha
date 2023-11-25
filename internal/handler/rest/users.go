package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	types2 "mocha/internal/types"
)

type UserHandler interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
}

func (s *restServer) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		// channelIDStr이 올바른 int64로 변환되지 않은 경우 에러 처리
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user Id"})
		return
	}

	user, err := s.rdb.GetUser(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 레코드를 찾지 못한 경우 404 에러 반환
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			return
		}
		// 다른 에러가 발생한 경우 500 에러 반환
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get user"})
		return
	}

	json.NewEncoder(w).Encode(map[string]types2.User{"user": *user})
}

func (s *restServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cu types2.CreateUser
	err := json.NewDecoder(r.Body).Decode(&cu)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	// validator를 사용하여 필수 파라미터 체크
	validate := validator.New()
	if err := validate.Struct(cu); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	foundUser, err := s.rdb.GetUserByEmail(cu.Email)
	if foundUser != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Duplicated user"})
		return
	}

	var user = types2.User{
		Name:      cu.Name,
		Password:  cu.Password,
		Email:     cu.Email,
		Age:       cu.Age,
		Gender:    cu.Gender,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	createdUser, err := s.rdb.CreateUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user"})
		return
	}

	// create bot chat
	/*
		1. bot user 가져오기 이름 meow
		2. bot user의 user id로 chat방 만들기
	*/

	botUser, err := s.rdb.GetBotByName("meow")
	var channel = types2.Channel{
		Type:            "bot",
		Name:            "meow-meow",
		HostUserId:      user.Id,
		LastMessageId:   0,
		LastDecryptedId: 0,
		CreatedAt:       time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:       time.Now().UTC().Format(time.RFC3339),
	}
	createdBotChannel, err := s.rdb.CreateChannel(&channel)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create bot channel"})
		return
	}

	// channel user 생성
	var cuser = types2.ChannelUser{
		ChannelId:         createdBotChannel.Id,
		UserId:            user.Id,
		LastSeenMessageId: 0,
	}
	err = s.rdb.CreateChannelUser(&cuser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create channel_user"})
		return
	}

	cuser = types2.ChannelUser{
		ChannelId:         createdBotChannel.Id,
		UserId:            botUser.Id,
		LastSeenMessageId: 0,
	}
	err = s.rdb.CreateChannelUser(&cuser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create channel_user"})
		return
	}

	json.NewEncoder(w).Encode(map[string]types2.User{"user": *createdUser})
}

func (s *restServer) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var lu types2.LoginUser
	err := json.NewDecoder(r.Body).Decode(&lu)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	// validator를 사용하여 필수 파라미터 체크
	validate := validator.New()
	if err := validate.Struct(lu); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	foundUser, err := s.rdb.GetUserByEmail(lu.Email)
	if s.handleError(w, err) {
		return
	}
	if foundUser == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Duplicated user"})
		return
	}

	if foundUser.Password != lu.Password {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid Password"})
		return
	}

	json.NewEncoder(w).Encode(map[string]types2.User{"user": *foundUser})
}

func (s *restServer) handleError(w http.ResponseWriter, err error) bool {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 레코드를 찾지 못한 경우 404 에러 반환
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			return true
		}
		// 다른 에러가 발생한 경우 500 에러 반환
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get user"})
		return true
	}
	return false
}
