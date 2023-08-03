package types

// Conversation struct는 채팅방의 데이터 구조를 나타냅니다.
type CreateConversation struct {
	Name        string  `json:"name,omitempty"`
	HostUserId  int64   `json:"host_user_id"`
	JoinedUsers []int64 `json:"joined_users"`
}

type Conversation struct {
	Id              int64  `gorm:"primaryKey;autoIncrement"`
	Type            string `gorm:"not null"`
	Name            string `gorm:"not null"`
	HostUserId      int64  `gorm:"not null"`
	LastMessageId   int64
	LastDecryptedId int64
	CreatedAt       string
	UpdatedAt       string
}

type ConversationUser struct {
	ConversationId    int64 `gorm:"primaryKey"`
	UserId            int64 `gorm:"primaryKey"`
	LastSeenMessageId int64
}
