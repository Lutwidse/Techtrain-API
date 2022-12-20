package service

import (
	"fmt"
	"math/rand"
	"net/http"
	"sort"

	"github.com/Lutwidse/Techtrain-API/internal/model/data"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GachaService struct {
	Db         *gorm.DB
	Gacha      data.Gacha
	GachaArray data.GachaArray
}

type GachaRequest struct {
	Times int `json:"times"`
}

type GachaResponse struct {
	CharacterId int    `json:"CharacterID"`
	Name        string `json:"name"`
}

type IndexChunk struct {
	From, To int
}

func IndexChunks(length int, chunkSize int) <-chan IndexChunk {
	ch := make(chan IndexChunk)

	go func() {
		defer close(ch)

		for i := 0; i < length; i += chunkSize {
			idx := IndexChunk{i, i + chunkSize}
			if length < idx.To {
				idx.To = length
			}
			ch <- idx
		}
	}()

	return ch
}

func (s *GachaService) Draw(c *gin.Context) {
	var gachaRequest GachaRequest
	var gachaResponse []GachaResponse

	if err := c.BindJSON(&gachaRequest); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Draw Times Required"})
		return
	}

	token := c.GetHeader("x-token")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Token Required"})
		return
	}

	drawTimes := gachaRequest.Times
	if drawTimes <= 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Draw Times Invalid"})
		return
	}

	result := s.Db.Table("gachas").Find(&s.GachaArray)
	if result.RowsAffected == 0 {
		var dummy [0]int
		c.JSON(http.StatusOK, gin.H{"results": dummy})
		return
	}

	gachaWeights := make([]int, 0)
	var gachaNames = make(map[int]string)
	for i := 0; i < len(s.GachaArray); i++ {
		weight := s.GachaArray[i].Weight
		gachaWeights = append(gachaWeights, weight)
		gachaNames[i] = s.GachaArray[i].Name
	}

	boundaries := make([]int, len(gachaWeights)+1)
	for i := 1; i < len(boundaries); i++ {
		boundaries[i] = boundaries[i-1] + gachaWeights[i-1]
	}

	// ガチャ
	boundaryLast := boundaries[len(boundaries)-1]
	counter := make([]int, len(gachaWeights))
	gachaDraws := make([]int, 0)
	for i := 0; i < drawTimes; i++ {
		x := rand.Intn(boundaryLast) + 1
		idx := sort.SearchInts(boundaries, x) - 1
		counter[idx]++
		gachaDraws = append(gachaDraws, idx+1)
	}

	// ガチャ結果 - 確認
	for i := 0; i < len(gachaWeights); i++ {
		fmt.Printf(
			"|%d|%f|%f|\n",
			gachaWeights[i],
			float64(gachaWeights[i])/float64(boundaryLast),
			float64(counter[i])/float64(drawTimes),
		)
	}

	characterBatch := make([]data.Character, 0)
	for i := 0; i < int(len(gachaDraws)); i++ {
		characterID := gachaDraws[i] - 1
		gachaName := gachaNames[characterID]

		characterBatch = append(characterBatch, data.Character{CharacterId: characterID, XToken: token})
		gachaResponse = append(gachaResponse, GachaResponse{CharacterId: characterID, Name: gachaName})
	}

	if len(characterBatch) >= 10000 {
		for idx := range IndexChunks(len(characterBatch), 10000) {
			s.Db.Table("characters").Create(characterBatch[idx.From:idx.To])
		}
	} else {
		s.Db.Table("characters").Create(&characterBatch)
	}

	if len(gachaResponse) > 100 {
		dummy := gachaResponse[:100]
		c.JSON(http.StatusOK, gin.H{"results": dummy})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": gachaResponse})
}
