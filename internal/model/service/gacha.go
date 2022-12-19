package service

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type GachaService struct {
	db    *gorm.DB
	gacha Gacha
}

type GachaRequest struct {
	Times int `json:"times"`
}

type GachaResponse struct {
	CharacterId int    `json:"CharacterID"`
	Name        string `json:"name"`
}

func (s *GachaService) Draw(c *gin.Context) {
	var gacha Gacha
	var gachaReq GachaRequest
	var gachaRes []GachaResponse

	if err := c.BindJSON(&gachaReq); err != nil {
		log.Fatal(err)
	}
	token := c.GetHeader("x-token")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Token Required"})
		return
	}

	gachaDraw := gachaReq
	times := gachaDraw.Times

	for i := 0; i < times; i++ {
		rnd := rand.Intn(100)
		result := s.db.Where("min_ratio <= ?", rnd).Where("max_ratio >= ?", rnd).First(&gacha)
		characterId := result.Value.(*Gacha).CharacterId
		name := result.Value.(*Gacha).Name

		s.db.Exec("INSERT INTO `techtrain_db`.`characters` (`name`, `character_id`, `x_token`) VALUES (?, ?, ?)", name, characterId, token)
		gachaRes = append(gachaRes, GachaResponse{CharacterId: characterId, Name: name})
	}
	c.JSON(http.StatusOK, gin.H{"results": gachaRes})
}
