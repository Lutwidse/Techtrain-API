package data

type User struct {
	Name   string `gorm:"column:name"`
	XToken string `gorm:"column:x_token"`
}