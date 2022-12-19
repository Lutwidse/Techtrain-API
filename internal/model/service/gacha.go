package service

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/Lutwidse/Techtrain-API/internal/model/data"
)

type GachaService struct {
	Db    *gorm.DB
	Gacha data.Gacha
}

type GachaRequest struct {
	Times int `json:"times"`
}

type GachaResponse struct {
	CharacterId int    `json:"CharacterID"`
	Name        string `json:"name"`
}

func (s *GachaService) Draw(c *gin.Context) {
	var gacha data.Gacha
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
		result := s.Db.Where("min_ratio <= ?", rnd).Where("max_ratio >= ?", rnd).First(&gacha)
		characterId := result.Value.(*data.Gacha).CharacterId
		name := result.Value.(*data.Gacha).Name

		s.Db.Exec("INSERT INTO `techtrain_db`.`characters` (`name`, `character_id`, `x_token`) VALUES (?, ?, ?)", name, characterId, token)
		gachaRes = append(gachaRes, GachaResponse{CharacterId: characterId, Name: name})
	}
	c.JSON(http.StatusOK, gin.H{"results": gachaRes})
}
