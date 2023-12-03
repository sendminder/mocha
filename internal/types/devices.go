package types

type CreateDevice struct {
	UserId    int64  `json:"user_id" validate:"required"`
	PushToken string `json:"push_token" validate:"required"`
	Platform  string `json:"platform" validate:"required"`
	Version   string `json:"version"`
}

type Device struct {
	Id        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    int64  `gorm:"primaryKey" json:"user_id"`
	PushToken string `gorm:"not null" json:"push_token"`
	Platform  string `gorm:"not null" json:"platform"`
	Version   string `json:"version"`
	Activated bool   `gorm:"not null" json:"activated"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
