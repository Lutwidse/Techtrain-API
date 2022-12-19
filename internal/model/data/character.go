package data

type Character []struct {
	Name        string `gorm:"column:name"`
	CharacterId int    `gorm:"column:character_id"`
	xToken      string `gorm:"column:x_token"`
}