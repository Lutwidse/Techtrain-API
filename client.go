package techtrain_api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
)

type TechtrainClient struct {
	user      *User
	gacha     *Gacha
	character *Character
}

func NewTechtrainClient() *TechtrainClient {
	return &TechtrainClient{user: &User{}, gacha: &Gacha{}, character: &Character{}}
}

func (client *TechtrainClient) Server() {
	router := gin.Default()

	userApi := router.Group("user")
	{
		userApi.POST("/create", client.user.UserCreate)
		userApi.GET("/get", client.user.UserGet)
		userApi.PUT("/update", client.user.UserUpdate)
	}

	gachaApi := router.Group("gacha")
	{
		gachaApi.POST("/draw", client.gacha.GachaDraw)
	}

	characterApi := router.Group("character")
	{
		characterApi.GET("/list", client.character.CharacterList)
	}

	c := cors.AllowAll()
	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", handler))
}
