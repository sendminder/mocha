package db

import (
	"math/rand"

	"mocha/internal/types"
)

type BotRecorder interface {
	CreateBot(bot *types.Bot) (*types.Bot, error)
	GetBotByName(name string) (*types.Bot, error)
}

func (db *rdb) CreateBot(bot *types.Bot) (*types.Bot, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	result := db.con[randIdx].Create(bot)
	if result.Error != nil {
		return nil, result.Error
	}
	return bot, nil
}

func (db *rdb) GetBotByName(name string) (*types.Bot, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	var bot types.Bot
	result := db.con[randIdx].First(&bot, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return &bot, nil
}
