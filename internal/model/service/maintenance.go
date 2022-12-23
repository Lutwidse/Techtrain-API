package service

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MaintenanceService struct {
	Db        *gorm.DB
	Srv       *http.Server
	Operation chan int
}

type MaintenanceRequest struct {
	Operation int `json:"operation"`
}

func (s *MaintenanceService) DebugSleep(c *gin.Context) {
	log.Println("Sleep Start...")
	time.Sleep(10 * time.Second)
	log.Println("Sleep End...")
	c.JSON(http.StatusInternalServerError, gin.H{"code": "Success"})
}

func (s *MaintenanceService) Fetch(c *gin.Context) {
	var maintenanceRequest MaintenanceRequest

	if err := c.BindJSON(&maintenanceRequest); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "No Operation Found"})
		return
	}

	// Maintenance Mode
	if maintenanceRequest.Operation > 0 && maintenanceRequest.Operation < 10 {
		s.Operation <- maintenanceRequest.Operation
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Operation"})
		return
	}
}
