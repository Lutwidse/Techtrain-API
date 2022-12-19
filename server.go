package techtrain_api

import (
	"log"
	"net/http"

	"/internal/model/service/"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
)

type TechtrainServer struct {
	user      *UserService
	gacha     *GachaService
	character *CharacterService
}

func NewTechtrainServer(db *gorm.DB) *TechtrainServer {
	return &TechtrainServer{user: &UserService{db, User{}}, gacha: &GachaService{db, Gacha{}}, character: &CharacterService{db, Character{}}}
}

func (server *TechtrainServer) Server() {
	router := gin.Default()

	userAPI := router.Group("user")
	{
		userAPI.POST("/create", server.user.Create)
		userAPI.GET("/get", server.user.Get)
		userAPI.PUT("/update", server.user.Update)
	}

	gachaAPI := router.Group("gacha")
	{
		gachaAPI.POST("/draw", server.gacha.Draw)
	}

	characterAPI := router.Group("character")
	{
		characterAPI.GET("/list", server.character.List)
	}

	c := cors.AllowAll()
	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", handler))
}
