package service

import (
	"net/http"

	"github.com/Lutwidse/Techtrain-API/internal/model/data"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const characterResponseLimit = 100

// CharacterService is Object
type CharacterService struct {
	Db             *gorm.DB
	Character      data.Character
	CharacterArray data.CharacterArray
}

// CharacterResponse is response struct of List
type CharacterResponse struct {
	Name        string `json:"name"`
	CharacterID int    `json:"characterID"`
}

// List returns characters owned by the user
func (s *CharacterService) List(c *gin.Context) {
	characterResponse := make([]CharacterResponse, 0)

	token := c.GetHeader("x-token")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Token Required"})
		return
	}

	characterResult := s.Db.Table("characters").Where("x_token = ?", token).Limit(characterResponseLimit).Find(&s.CharacterArray)
	if characterResult.RowsAffected == 0 {
		var dummy [0]int
		c.JSON(http.StatusOK, gin.H{"characters": dummy})
		return
	}

	var gachaArray data.GachaArray
	gachaResult := s.Db.Table("gachas").Find(&gachaArray)
	if gachaResult.RowsAffected == 0 {
		var dummy [0]int
		c.JSON(http.StatusOK, gin.H{"results": dummy})
		return
	}

	var gachaNames = make(map[int]string)
	for i := 0; i < len(gachaArray); i++ {
		gachaNames[i] = gachaArray[i].Name
	}

	for i := 0; i < int(characterResult.RowsAffected); i++ {
		characterID := s.CharacterArray[i].CharacterID

		characterResponse = append(characterResponse, CharacterResponse{CharacterID: characterID, Name: gachaNames[characterID]})
	}
	c.JSON(http.StatusOK, gin.H{"characters": characterResponse})
}
