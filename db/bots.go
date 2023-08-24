package db

import (
	"math/rand"
	"mocha/types"
)

func CreateBot(bot *types.Bot) (*types.Bot, error) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	result := dbConnections[randIdx].Create(bot)
	if result.Error != nil {
		return nil, result.Error
	}
	return bot, nil
}

func GetBotByName(name string) (*types.Bot, error) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	var bot types.Bot
	result := dbConnections[randIdx].First(&bot, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return &bot, nil
}
