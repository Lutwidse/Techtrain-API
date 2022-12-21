package techtrain_api

import (
	"log"
	"net/http"
	"sync"

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

func NewTechtrainServer(db *gorm.DB, wg *sync.WaitGroup) *TechtrainServer {
	return &TechtrainServer{
		user:        &service.UserService{db, wg, data.User{}},
		gacha:       &service.GachaService{db, wg, data.Gacha{}, data.GachaArray{}},
		character:   &service.CharacterService{db, wg, data.Character{}, data.CharacterArray{}},
		maintenance: &service.MaintenanceService{db, nil, wg}}
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

	log.Println(server.maintenance.Srv.ListenAndServe())
}
