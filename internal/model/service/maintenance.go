package service

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MaintenanceService is object
type MaintenanceService struct {
	Db        *gorm.DB
	Srv       *http.Server
	Operation chan int
}

// MaintenanceRequest is request struct of FetchOperation
type MaintenanceRequest struct {
	Operation int `json:"operation"`
}

// DebugSleep is just debug to check FetchOperation
func (s *MaintenanceService) DebugSleep(c *gin.Context) {
	log.Println("Sleep Start...")
	time.Sleep(10 * time.Second)
	log.Println("Sleep End...")
	c.JSON(http.StatusInternalServerError, gin.H{"code": "Success"})
}

// FetchOperation is receiver for FetchPoll
func (s *MaintenanceService) FetchOperation(c *gin.Context) {
	var maintenanceRequest MaintenanceRequest

	if err := c.BindJSON(&maintenanceRequest); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "No Operation Found"})
		return
	}

	s.Operation <- maintenanceRequest.Operation
}

// FetchPoll is polling for a wait operation
func (s *MaintenanceService) FetchPoll() {
	for {
		select {
		case op := <-s.Operation:
			switch op {
			case 1:
				log.Println("Shutdown...")

				ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
				defer cancel()
				if err := s.Srv.Shutdown(ctx); err != nil {
					log.Fatal("Force Shutdown", err)
				}
			}
		}
	}
}
