package service

import (
	"encoding/json"
	"errors"
	"mocha/db"
	"mocha/types"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		// conversationIDStr이 올바른 int64로 변환되지 않은 경우 에러 처리
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user Id"})
		return
	}

	user, err := db.GetUser(userId)
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

	json.NewEncoder(w).Encode(map[string]types.User{"user": *user})
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cu types.CreateUser
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

	foundUser, err := db.GetUserByEmail(cu.Email)
	if foundUser != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Duplicated user"})
		return
	}

	var user = types.User{
		Name:      cu.Name,
		Password:  cu.Password,
		Email:     cu.Email,
		Age:       cu.Age,
		Gender:    cu.Gender,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	createdUser, err := db.CreateUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user"})
		return
	}

	json.NewEncoder(w).Encode(map[string]types.User{"user": *createdUser})
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var lu types.LoginUser
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

	foundUser, err := db.GetUserByEmail(lu.Email)
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

	json.NewEncoder(w).Encode(map[string]types.User{"user": *foundUser})
}
