package data

type Character struct {
	CharacterId int    `gorm:"column:character_id"`
	XToken      string `gorm:"column:x_token"`
}

type CharacterArray []struct {
	CharacterId int    `gorm:"column:character_id"`
	XToken      string `gorm:"column:x_token"`
}