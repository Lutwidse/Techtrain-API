package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type UserService struct {
	db   *gorm.DB
	user User
}

func (s *UserService) Create(c *gin.Context) {
	if err := c.BindJSON(&s.user); err != nil {
		log.Fatal(err)
	}
	token := uuid.New().String()
	userReq := User{Name: s.user.Name, xToken: token}

	result := s.db.Exec("INSERT INTO `techtrain_db`.`users` (`name`, `x_token`) VALUES (?, ?)", userReq.Name, userReq.xToken)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Already Registered"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *UserService) Get(c *gin.Context) {
	token := c.GetHeader("x-token")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Token Required"})
		return
	}

	result := s.db.First(&s.user, "x_token = ?", token)
	c.JSON(http.StatusOK, gin.H{"name": result.Value.(*User).Name})
}

func (s *UserService) Update(c *gin.Context) {
	if err := c.BindJSON(&s.user); err != nil {
		log.Fatal(err)
	}
	token := c.GetHeader("x-token")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Token Required"})
		return
	}

	userReq := User{Name: s.user.Name, xToken: token}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/user/get", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("x-token", userReq.xToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&s.user)
	if err != nil {
		panic(err)
	}

	newName := userReq.Name
	oldName := s.user.Name

	s.db.Exec("UPDATE `techtrain_db`.`users` SET `name` = ? WHERE (`name` = ?) and (`x_token` = ?)", newName, oldName, token)
	c.JSON(http.StatusOK, nil)
}
