package db

import (
	"log/slog"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"mocha/internal/types"
)

type RelationalDatabase interface {
	BotRecorder
	ChannelRecorder
	DeviceRecorder
	UserRecorder
}

var _ RelationalDatabase = (*rdb)(nil)

type rdb struct {
	con            []*gorm.DB
	loc            []*sync.RWMutex
	numConnections int
}

func NewRelationalDatabase(wg *sync.WaitGroup, numConnections int) RelationalDatabase {
	defer wg.Done()
	// 데이터베이스 연결
	dsn := "user=test password=test dbname=mocha host=localhost port=5431 sslmode=disable"
	var err error
	var cons []*gorm.DB
	var locks []*sync.RWMutex

	for i := 0; i < numConnections; i++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			slog.Error("Failed to connect to database", "error", err)
			return nil
		}
		cons = append(cons, db)
		locks = append(locks, new(sync.RWMutex))
		slog.Info("connected with db", "index", i)
	}

	// 테이블 생성
	err = cons[0].AutoMigrate(&types.Channel{}, &types.User{}, &types.ChannelUser{}, &types.Device{}, &types.Bot{})
	if err != nil {
		slog.Error("Failed to create table", "error", err)
		return nil
	}

	slog.Info("Table created successfully.")

	// email 컬럼에 인덱스를 추가하는 SQL 쿼리
	// 인덱스를 추가하려면 이 SQL 쿼리를 실행해야합니다.
	// dbConnections[0].Exec("CREATE INDEX idx_users_email ON users(email)")
	// dbConnections[0].Exec("CREATE INDEX idx_devices_push_token ON devices(push_token)")

	return &rdb{
		con:            cons,
		loc:            locks,
		numConnections: numConnections,
	}
}
