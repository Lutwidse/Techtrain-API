package service

import (
	"net/http"

	"github.com/Lutwidse/Techtrain-API/internal/model/data"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CharacterService struct {
	Db             *gorm.DB
	Character      data.Character
	CharacterArray data.CharacterArray
}

type CharacterResponse struct {
	Name        string `json:"name"`
	CharacterId int    `json:"CharacterID"`
}

func (s *CharacterService) List(c *gin.Context) {
	var characterResponse []CharacterResponse

	token := c.GetHeader("x-token")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Token Required"})
		return
	}

	result := s.Db.Table("characters").Where("x_token = ?", token).Find(&s.CharacterArray)
	if result.RowsAffected == 0 {
		var dummy [0]int
		c.JSON(http.StatusOK, gin.H{"characters": dummy})
		return
	}
	for i := 0; i < int(result.RowsAffected); i++ {
		characterId := s.CharacterArray[i].CharacterId

		characterResponse = append(characterResponse, CharacterResponse{CharacterId: characterId})
	}
	c.JSON(http.StatusOK, gin.H{"characters": characterResponse})
}
