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

	v1 := router.Group("/api/v1/user")
	{
		v1.POST("/create", client.user.UserCreate)
		v1.GET("/get", client.user.UserGet)
		v1.PUT("/update", client.user.UserUpdate)
	}
	c := cors.AllowAll()
	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe("127.0.0.1:5001", handler))
}
