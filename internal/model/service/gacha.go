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

const sqlInsertLimit = 10000
const gachaResponseLimit = 100

// GachaService is Object
type GachaService struct {
	Db         *gorm.DB
	Gacha      data.Gacha
	GachaArray data.GachaArray
}

// GachaRequest is request struct of Draw
type GachaRequest struct {
	Times int `json:"times"`
}

// GachaResponse is response struct of Draw
type GachaResponse struct {
	CharacterID int    `json:"CharacterID"`
	Name        string `json:"name"`
}

// IndexChunk is struct of IndexChunks
type IndexChunk struct {
	From, To int
}

// IndexChunks is function for separate array with chunk size
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

// Draw gacha and return results
func (s *GachaService) Draw(c *gin.Context) {
	var gachaRequest GachaRequest
	gachaResponse := make([]GachaResponse, 0)

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

		characterBatch = append(characterBatch, data.Character{CharacterID: characterID, XToken: token})
		gachaResponse = append(gachaResponse, GachaResponse{CharacterID: characterID, Name: gachaName})
	}

	if len(characterBatch) >= sqlInsertLimit {
		for idx := range IndexChunks(len(characterBatch), sqlInsertLimit) {
			s.Db.Table("characters").Create(characterBatch[idx.From:idx.To])
		}
	} else {
		s.Db.Table("characters").Create(&characterBatch)
	}

	if len(gachaResponse) > gachaResponseLimit {
		dummy := gachaResponse[:gachaResponseLimit]
		c.JSON(http.StatusOK, gin.H{"results": dummy})
	} else {
		c.JSON(http.StatusOK, gin.H{"results": gachaResponse})
	}
}
