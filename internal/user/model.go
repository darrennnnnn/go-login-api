package user

import (
	"time"
)

type User struct {
	Id       string `gorm:"primaryKey"`
	Email    string `gorm:"uniqueIndex;not null;column:email"`
	Username string `gorm:"uniqueIndex;not null;column:username"`
	Password string `gorm:"not null;column:password"`
}

type AccessToken struct {
	ID        string    `gorm:"type:varchar(255);primaryKey"`
	UserID    string    `gorm:"type:varchar(255);not null;index"`
	Revoked   bool      `gorm:"not null;default:false"`
	ExpiresAt time.Time `gorm:"type:timestamptz;not null"`

	User User `gorm:"foreignKey:UserID;references:Id"`
}