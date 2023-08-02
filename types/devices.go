package types

type CreateDevice struct {
	UserId    int64  `json:"user_id" validate:"required"`
	PushToken string `json:"push_token" validate:"required"`
	Platform  string `json:"platform" validate:"required"`
	Version   string `json:"version"`
}

type Device struct {
	Id        int64  `gorm:"primaryKey;autoIncrement"`
	UserId    int64  `gorm:"primaryKey;"`
	PushToken string `gorm:"not null"`
	Platform  string `gorm:"not null"`
	Version   string
	Activated bool `gorm:"not null"`
	CreatedAt string
	UpdatedAt string
}
