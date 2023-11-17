package db

import (
	"math/rand"

	"mocha/internal/types"
)

type UserRecorder interface {
	CreateUser(user *types.User) (*types.User, error)
	GetUser(id int64) (*types.User, error)
	GetUserByEmail(email string) (*types.User, error)
	UpdateUser(id int64, updatedUser types.User) (*types.User, error)
	DeleteUser(id int64) error
}

func (db *rdb) CreateUser(user *types.User) (*types.User, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	result := db.con[randIdx].Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (db *rdb) GetUser(id int64) (*types.User, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	var user types.User
	result := db.con[randIdx].First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (db *rdb) GetUserByEmail(email string) (*types.User, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	var user types.User
	result := db.con[randIdx].Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (db *rdb) UpdateUser(id int64, updatedUser types.User) (*types.User, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	result := db.con[randIdx].Model(&types.User{}).Where("Id = ?", id).Updates(updatedUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &updatedUser, nil
}

func (db *rdb) DeleteUser(id int64) error {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	var user types.User
	result := db.con[randIdx].First(&user, id)
	if result.Error != nil {
		return result.Error
	}
	result = db.con[randIdx].Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
