package user

type User struct {
	Id       string `gorm:"primaryKey"`
	Name     string `gorm:"not null;column:name"`
	UserName string `gorm:"uniqueIndex;not null;column:username"`
	Password string `gorm:"not null;column:password"`
}