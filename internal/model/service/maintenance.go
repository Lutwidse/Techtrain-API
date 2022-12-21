package service

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MaintenanceService struct {
	Db  *gorm.DB
	Srv *http.Server
	Wg  *sync.WaitGroup
}

type MaintenanceRequest struct {
	Operation int `json:"operation"`
}

func (s *MaintenanceService) DebugSleep(c *gin.Context) {
	s.Wg.Add(1)
	log.Println("Sleep Start...")
	time.Sleep(10 * time.Second)
	log.Println("Sleep End...")
	c.JSON(http.StatusInternalServerError, gin.H{"code": "Success"})
	s.Wg.Done()
}

func (s *MaintenanceService) Fetch(c *gin.Context) {
	var maintenanceRequest MaintenanceRequest

	if err := c.BindJSON(&maintenanceRequest); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "No Operation Found"})
		return
	}
	// Maintenance Mode
	if maintenanceRequest.Operation == 1 {
		s.Srv.Handler = nil
		s.Wg.Wait()
		log.Println("Shutting down server...")

		s.Srv.Shutdown(c)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Operation"})
		return
	}
}
