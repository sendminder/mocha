package types

type LoginUser struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

type CreateUser struct {
	Name     string `json:"name,omitempty" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
}

type User struct {
	Id        int64  `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"not null"`
	Password  string `gorm:"not null"`
	Email     string `gorm:"not null"`
	Age       int
	Gender    string
	CreatedAt string
	UpdatedAt string
}
