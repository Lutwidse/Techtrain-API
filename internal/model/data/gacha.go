package data

type Gacha struct {
	CharacterId int    `gorm:"column:character_id"`
	Name        string `gorm:"column:name"`
	MinRatio    int    `gorm:"column:min_ratio"`
	MaxRatio    int    `gorm:"column:max_ratio"`
}