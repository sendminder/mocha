package db

import (
	"math/rand"
	"mocha/types"
)

func CreateDevice(device *types.Device) (*types.Device, error) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	result := dbConnections[randIdx].Create(device)
	if result.Error != nil {
		return nil, result.Error
	}
	return device, nil
}

func GetDevice(id int64) (*types.Device, error) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	var device types.Device
	result := dbConnections[randIdx].First(&device, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &device, nil
}

func GetDeviceByUserId(userId int64) (*types.Device, error) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	var device types.Device
	result := dbConnections[randIdx].Where("user_id = ?", userId).First(&device)
	if result.Error != nil {
		return nil, result.Error
	}
	return &device, nil
}

func GetDeviceByPushToken(pushToken string) (*types.Device, error) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	var device types.Device
	result := dbConnections[randIdx].Where("push_token = ?", pushToken).First(&device)
	if result.Error != nil {
		return nil, result.Error
	}
	return &device, nil
}

func UpdateDevice(id int64, updatedUser types.Device) (*types.Device, error) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	result := dbConnections[randIdx].Model(&types.Device{}).Where("Id = ?", id).Updates(updatedUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &updatedUser, nil
}

func DeleteDevice(id int64) error {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	var device types.Device
	result := dbConnections[randIdx].First(&device, id)
	if result.Error != nil {
		return result.Error
	}
	result = dbConnections[randIdx].Delete(&device)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
