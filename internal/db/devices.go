package db

import (
	"math/rand"

	"mocha/internal/types"
)

type DeviceRecorder interface {
	CreateDevice(device *types.Device) (*types.Device, error)
	GetDevice(id int64) (*types.Device, error)
	GetDeviceByUserID(userID int64) (*types.Device, error)
	GetDeviceByPushToken(pushToken string) (*types.Device, error)
	UpdateDevice(id int64, updatedUser types.Device) (*types.Device, error)
	DeleteDevice(id int64) error
}

func (db *rdb) CreateDevice(device *types.Device) (*types.Device, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	result := db.con[randIdx].Create(device)
	if result.Error != nil {
		return nil, result.Error
	}
	return device, nil
}

func (db *rdb) GetDevice(id int64) (*types.Device, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	var device types.Device
	result := db.con[randIdx].First(&device, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &device, nil
}

func (db *rdb) GetDeviceByUserID(userID int64) (*types.Device, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	var device types.Device
	result := db.con[randIdx].Where("user_id = ?", userID).First(&device)
	if result.Error != nil {
		return nil, result.Error
	}
	return &device, nil
}

func (db *rdb) GetDeviceByPushToken(pushToken string) (*types.Device, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	var device types.Device
	result := db.con[randIdx].Where("push_token = ?", pushToken).First(&device)
	if result.Error != nil {
		return nil, result.Error
	}
	return &device, nil
}

func (db *rdb) UpdateDevice(id int64, updatedUser types.Device) (*types.Device, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	result := db.con[randIdx].Model(&types.Device{}).Where("Id = ?", id).Updates(updatedUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &updatedUser, nil
}

func (db *rdb) DeleteDevice(id int64) error {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	var device types.Device
	result := db.con[randIdx].First(&device, id)
	if result.Error != nil {
		return result.Error
	}
	result = db.con[randIdx].Delete(&device)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
