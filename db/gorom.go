package db

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbConnections []*gorm.DB
var dbLocks []*sync.RWMutex

const numConnections = 10

type Conversation struct {
	ID            int64  `gorm:"primaryKey;autoIncrement"`
	Type          string `gorm:"not null"`
	Name          string `gorm:"not null"`
	HostUserID    int64  `gorm:"not null"`
	LastMessageID int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type User struct {
	ID        int64  `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Password  string `gorm:"not null"`
	Email     string `gorm:"not null"`
	Age       int    `gorm:"not null"`
	Gender    string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ConversationUser struct {
	ConversationID    int64 `gorm:"primaryKey"`
	UserID            int64 `gorm:"primaryKey"`
	LastSeenMessageID int64
}

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
	err = dbConnections[0].AutoMigrate(&Conversation{}, &User{}, &ConversationUser{})
	if err != nil {
		fmt.Println("Failed to create table:", err)
		return
	}

	fmt.Println("Table created successfully.")
}

func CreateConversation(conversation *Conversation) (*Conversation, error) {
	radIdx := rand.Intn(10)
	dbLocks[radIdx].Lock()
	defer dbLocks[radIdx].Unlock()

	result := dbConnections[radIdx].Create(conversation)
	if result.Error != nil {
		return nil, result.Error
	}
	return conversation, nil
}

func GetConversationByID(conversationID int64) (*Conversation, error) {
	radIdx := rand.Intn(10)
	dbLocks[radIdx].Lock()
	defer dbLocks[radIdx].Unlock()

	var conversation Conversation
	result := dbConnections[radIdx].First(&conversation, conversationID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &conversation, nil
}

func GetUserConversations(uesrId int64) ([]Conversation, error) {
	radIdx := rand.Intn(10)
	dbLocks[radIdx].Lock()
	defer dbLocks[radIdx].Unlock()

	// GORM을 이용하여 conversations 테이블과 conversation_users 테이블을 JOIN하여 특정 사용자가 참여한 채팅방들을 가져옴
	var conversations []Conversation
	result := dbConnections[radIdx].
		Joins("JOIN conversation_users ON conversations.id = conversation_users.conversation_id").
		Where("conversation_users.user_id = ?", uesrId).
		Find(&conversations)

	if result.Error != nil {
		return nil, result.Error
	}

	return conversations, nil
}

func UpdateConversation(conversation *Conversation) error {
	radIdx := rand.Intn(10)
	dbLocks[radIdx].Lock()
	defer dbLocks[radIdx].Unlock()

	result := dbConnections[radIdx].Save(conversation)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteConversation(conversationID uint) error {
	radIdx := rand.Intn(10)
	dbLocks[radIdx].Lock()
	defer dbLocks[radIdx].Unlock()

	result := dbConnections[radIdx].Delete(&Conversation{}, conversationID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func CreateConversationUser(cuser *ConversationUser) error {
	radIdx := rand.Intn(10)
	dbLocks[radIdx].Lock()
	defer dbLocks[radIdx].Unlock()

	result := dbConnections[radIdx].Create(cuser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
