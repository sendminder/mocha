package types

type User struct {
	Id        int64  `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"not null"`
	Password  string `gorm:"not null"`
	Email     string `gorm:"not null"`
	Age       int    `gorm:"not null"`
	Gender    string `gorm:"not null"`
	CreatedAt string
	UpdatedAt string
}
