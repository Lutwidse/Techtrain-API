package techtrain_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Character []struct {
	Name        string `gorm:"column:name"`
	CharacterId int    `gorm:"column:character_id"`
	xToken      string `gorm:"column:x_token"`
}
type CharacterService struct {
	db        *gorm.DB
	character Character
}

type CharacterResponse struct {
	Name        string `json:"name"`
	CharacterId int    `json:"CharacterID"`
}

func (s *CharacterService) List(c *gin.Context) {
	var charaRes []CharacterResponse

	token := c.GetHeader("x-token")
	result := s.db.Table("characters").Where("x_token = ?", token).Find(&s.character)
	for i := 0; i < int(result.RowsAffected); i++ {
		name := s.character[i].Name
		characterId := s.character[i].CharacterId

		charaRes = append(charaRes, CharacterResponse{Name: name, CharacterId: characterId})
	}
	c.JSON(http.StatusOK, gin.H{"characters": charaRes})
}
