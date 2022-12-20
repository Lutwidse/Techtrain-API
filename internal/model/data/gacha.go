package data

type Gacha struct {
	CharacterId int    `gorm:"column:character_id"`
	Weight      int    `gorm:"column:weight"`
	Name        string `gorm:"column:name"`
}
type GachaArray []struct {
	CharacterId int    `gorm:"column:character_id"`
	Weight      int    `gorm:"column:weight"`
	Name        string `gorm:"column:name"`
}
