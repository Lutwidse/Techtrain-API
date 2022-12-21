package data

type Character struct {
	CharacterID int    `gorm:"column:character_id"`
	XToken      string `gorm:"column:x_token"`
}

type CharacterArray []struct {
	CharacterID int    `gorm:"column:character_id"`
	XToken      string `gorm:"column:x_token"`
}