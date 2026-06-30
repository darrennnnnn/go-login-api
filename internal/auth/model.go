package auth

import (
	"time"

	"github.com/darrennnnnn/go-login-api/internal/user"
)

type AccessToken struct {
	ID        string    `gorm:"type:varchar(255);primaryKey"`
	UserID    string    `gorm:"type:varchar(255);not null;index"`
	Revoked   bool      `gorm:"not null;default:false"`
	ExpiresAt time.Time `gorm:"type:timestamptz;not null"`

	User user.User `gorm:"foreignKey:UserID;references:ID"`
}
