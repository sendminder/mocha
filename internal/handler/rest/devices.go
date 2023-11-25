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
	"mocha/internal/types"
)

type DeviceHandler interface {
	GetDevice(w http.ResponseWriter, r *http.Request)
	CreateDevice(w http.ResponseWriter, r *http.Request)
}

func (s *restServer) GetDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	deviceId, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		// channelIDStr이 올바른 int64로 변환되지 않은 경우 에러 처리
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid device Id"})
		return
	}

	device, err := s.rdb.GetDevice(deviceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 레코드를 찾지 못한 경우 404 에러 반환
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Device not found"})
			return
		}
		// 다른 에러가 발생한 경우 500 에러 반환
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get device"})
		return
	}

	json.NewEncoder(w).Encode(map[string]types.Device{"device": *device})
}

func (s *restServer) CreateDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cd types.CreateDevice
	err := json.NewDecoder(r.Body).Decode(&cd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	// validator를 사용하여 필수 파라미터 체크
	validate := validator.New()
	if err := validate.Struct(cd); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	foundDevice, err := s.rdb.GetDeviceByPushToken(cd.PushToken)
	if foundDevice != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Duplicated device"})
		return
	}

	var device = types.Device{
		UserId:    cd.UserId,
		PushToken: cd.PushToken,
		Platform:  cd.Platform,
		Version:   cd.Version,
		Activated: true,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	createdDevice, err := s.rdb.CreateDevice(&device)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create device"})
		return
	}

	json.NewEncoder(w).Encode(map[string]types.Device{"device": *createdDevice})
}
