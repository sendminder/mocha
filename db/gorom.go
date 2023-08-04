package db

import (
	"log"
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
			log.Println("Failed to connect to database:", err)
			return
		}
		dbConnections = append(dbConnections, db)
		dbLocks = append(dbLocks, new(sync.RWMutex))
		log.Println("connected with db", i)
	}

	// 테이블 생성
	err = dbConnections[0].AutoMigrate(&types.Conversation{}, &types.User{}, &types.ConversationUser{}, &types.Device{})
	if err != nil {
		log.Println("Failed to create table:", err)
		return
	}

	log.Println("Table created successfully.")

	// email 컬럼에 인덱스를 추가하는 SQL 쿼리
	// 인덱스를 추가하려면 이 SQL 쿼리를 실행해야합니다.
	// dbConnections[0].Exec("CREATE INDEX idx_users_email ON users(email)")
	// dbConnections[0].Exec("CREATE INDEX idx_devices_push_token ON devices(push_token)")
}
