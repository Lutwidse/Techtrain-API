package data

type User struct {
	Name   string `gorm:"column:name"`
	xToken string `gorm:"column:x_token"`
}