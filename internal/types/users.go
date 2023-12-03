package types

type LoginUser struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email"    validate:"required"`
}

type CreateUser struct {
	Name     string `json:"name,omitempty" validate:"required"`
	Password string `json:"password"       validate:"required"`
	Email    string `json:"email"          validate:"required"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
}

type User struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"not null"                 json:"name"`
	Password  string `gorm:"not null"                 json:"password"`
	Email     string `gorm:"not null"                 json:"email"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	IsBot     bool   `gorm:"default:false"            json:"is_bot"`
}

type Bot struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64  `gorm:"not null"                 json:"user_id"`
	Name      string `gorm:"not null"                 json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
