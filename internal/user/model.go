package user

type User struct {
	Id       string `gorm:"primaryKey"`
	Email    string `gorm:"uniqueIndex;not null;column:email"`
	Username string `gorm:"uniqueIndex;not null;column:username"`
	Password string `gorm:"not null;column:password"`
}