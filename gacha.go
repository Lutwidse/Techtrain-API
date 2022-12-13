package techtrain_api

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Gacha struct {
	CharacterId int    `gorm:"column:character_id"`
	Name        string `gorm:"column:name"`
	MinRatio    int    `gorm:"column:min_ratio"`
	MaxRatio    int    `gorm:"column:max_ratio"`
}

type GachaRequest struct {
	Times int `json:"times"`
}

type GachaResponse struct {
	CharacterId int    `json:"CharacterID"`
	Name        string `json:"name"`
}

func (u *Gacha) GachaDraw(c *gin.Context) {
	db, err := gorm.Open("mysql", SqlArgs)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.LogMode(true)

	var gacha Gacha
	var gachaReq GachaRequest
	var gachaRes []GachaResponse

	if err := c.BindJSON(&gachaReq); err != nil {
		log.Fatal(err)
	}
	token := c.GetHeader("x-token")

	gachaDraw := gachaReq
	times := gachaDraw.Times

	for i := 0; i < times; i++ {
		rnd := rand.Intn(100)
		result := db.Where("min_ratio <= ?", rnd).Where("max_ratio >= ?", rnd).First(&gacha)
		characterId := result.Value.(*Gacha).CharacterId
		name := result.Value.(*Gacha).Name

		db.Exec("INSERT INTO `techtrain_db`.`characters` (`name`, `character_id`, `x_token`) VALUES (?, ?, ?)", name, characterId, token)
		gachaRes = append(gachaRes, GachaResponse{CharacterId: characterId, Name: name})
	}
	c.JSON(http.StatusOK, gin.H{"results": gachaRes})
}
