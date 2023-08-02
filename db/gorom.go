package db

import (
	"fmt"
	"mocha/types"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbConnections []*gorm.DB
var dbLocks []*sync.RWMutex

const numConnections = 10

func ConnectGorm(wg *sync.WaitGroup) {
	defer wg.Done()
	// 데이터베이스 연결
	dsn := "user=test password=test dbname=mocha host=localhost port=5431 sslmode=disable"
	var err error

	for i := 0; i < numConnections; i++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println("Failed to connect to database:", err)
			return
		}
		dbConnections = append(dbConnections, db)
		dbLocks = append(dbLocks, new(sync.RWMutex))
		fmt.Println("connected with db", i)
	}

	// 테이블 생성
	err = dbConnections[0].AutoMigrate(&types.Conversation{}, &types.User{}, &types.ConversationUser{})
	if err != nil {
		fmt.Println("Failed to create table:", err)
		return
	}

	fmt.Println("Table created successfully.")
}
