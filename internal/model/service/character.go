package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/Lutwidse/Techtrain-API/internal/model/data"
)

type CharacterService struct {
	Db        *gorm.DB
	Character data.Character
}

type CharacterResponse struct {
	Name        string `json:"name"`
	CharacterId int    `json:"CharacterID"`
}

func (s *CharacterService) List(c *gin.Context) {
	var charaRes []CharacterResponse

	token := c.GetHeader("x-token")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Token Required"})
		return
	}
	
	result := s.Db.Table("characters").Where("x_token = ?", token).Find(&s.Character)
	for i := 0; i < int(result.RowsAffected); i++ {
		name := s.Character[i].Name
		characterId := s.Character[i].CharacterId

		charaRes = append(charaRes, CharacterResponse{Name: name, CharacterId: characterId})
	}
	c.JSON(http.StatusOK, gin.H{"characters": charaRes})
}
