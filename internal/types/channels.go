package types

// Channel struct는 채팅방의 데이터 구조를 나타냅니다.
type CreateChannel struct {
	Name        string  `json:"name,omitempty"`
	HostUserID  int64   `json:"host_user_id"`
	JoinedUsers []int64 `json:"joined_users"`
}

type Channel struct {
	ID              int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Type            string `gorm:"not null"                 json:"type"`
	Name            string `gorm:"not null"                 json:"name"`
	HostUserID      int64  `gorm:"not null"                 json:"host_user_id"`
	LastMessageID   int64  `json:"last_message_id"`
	LastDecryptedID int64  `json:"last_decrypted_id"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type ChannelUser struct {
	ChannelID         int64 `gorm:"primaryKey"           json:"channel_id"`
	UserID            int64 `gorm:"primaryKey"           json:"user_id"`
	LastSeenMessageID int64 `json:"last_seen_message_id"`
}
