package db

import (
	"math/rand"
	"mocha/types"
)

func CreateUser(user *types.User) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	dbConnections[randIdx].Create(user)
}

func ReadUser(id int64) types.User {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	var user types.User
	dbConnections[randIdx].First(&user, id)
	return user
}

func UpdateUser(id int64, updatedUser types.User) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	dbConnections[randIdx].Model(&types.User{}).Where("Id = ?", id).Updates(updatedUser)
}

func DeleteUser(id int64) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	var user types.User
	dbConnections[randIdx].First(&user, id)
	dbConnections[randIdx].Delete(&user)
}
