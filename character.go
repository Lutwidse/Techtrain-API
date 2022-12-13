package techtrain_api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Character []struct {
	Name        string `gorm:"column:name"`
	CharacterId int    `gorm:"column:character_id"`
	xToken      string `gorm:"column:x_token"`
}

type CharacterResponse struct {
	Name        string `json:"name"`
	CharacterId int    `json:"CharacterID"`
}

func (u *Character) CharacterList(c *gin.Context) {
	db, err := gorm.Open("mysql", SqlArgs)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.LogMode(true)

	var chara Character
	var charaRes []CharacterResponse

	token := c.GetHeader("x-token")
	result := db.Table("characters").Where("x_token = ?", token).Find(&chara)
	for i := 0; i < int(result.RowsAffected); i++ {
		name := chara[i].Name
		characterId := chara[i].CharacterId

		charaRes = append(charaRes, CharacterResponse{Name: name, CharacterId: characterId})
	}
	c.JSON(http.StatusOK, gin.H{"characters": charaRes})
}
