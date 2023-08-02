package db

import (
	"math/rand"
	"mocha/types"
)

func CreateUser(user *types.User) (*types.User, error) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	result := dbConnections[randIdx].Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func GetUser(id int64) (*types.User, error) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	var user types.User
	result := dbConnections[randIdx].First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByEmail(email string) (*types.User, error) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	var user types.User
	result := dbConnections[randIdx].Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func UpdateUser(id int64, updatedUser types.User) (*types.User, error) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	result := dbConnections[randIdx].Model(&types.User{}).Where("Id = ?", id).Updates(updatedUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &updatedUser, nil
}

func DeleteUser(id int64) error {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	var user types.User
	result := dbConnections[randIdx].First(&user, id)
	if result.Error != nil {
		return result.Error
	}
	result = dbConnections[randIdx].Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
