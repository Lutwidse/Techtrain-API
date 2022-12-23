package techtrain_api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Lutwidse/Techtrain-API/internal/model/data"
	"github.com/Lutwidse/Techtrain-API/internal/model/service"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

type TechtrainServer struct {
	user        *service.UserService
	gacha       *service.GachaService
	character   *service.CharacterService
	maintenance *service.MaintenanceService
}

func NewTechtrainServer(db *gorm.DB) *TechtrainServer {
	return &TechtrainServer{
		user:        &service.UserService{db, data.User{}},
		gacha:       &service.GachaService{db, data.Gacha{}, data.GachaArray{}},
		character:   &service.CharacterService{db, data.Character{}, data.CharacterArray{}},
		maintenance: &service.MaintenanceService{db, nil, nil}}
}

func (server *TechtrainServer) Start() {
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

	maintenanceAPI := router.Group("maintenance")
	{
		maintenanceAPI.POST("/operation", server.maintenance.Fetch)
		maintenanceAPI.GET("/debugsleep", server.maintenance.DebugSleep)
	}

	c := cors.AllowAll()
	handler := c.Handler(router)

	server.maintenance.Srv = &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	operation := make(chan int, 1)
	server.maintenance.Operation = operation

	go func() {
		log.Println(server.maintenance.Srv.ListenAndServe())
	}()

	<-server.maintenance.Operation
	log.Println("Shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	if err := server.maintenance.Srv.Shutdown(ctx); err != nil {
		log.Fatal("Force Shutdown", err)
	}
}
