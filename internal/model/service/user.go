package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Lutwidse/Techtrain-API/internal/model/data"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type UserService struct {
	Db   *gorm.DB
	User data.User
}

func (s *UserService) Create(c *gin.Context) {
	if err := c.BindJSON(&s.User); err != nil {
		log.Fatal(err)
	}
	token := uuid.New().String()
	userReq := data.User{Name: s.User.Name, XToken: token}

	result := s.Db.Exec("INSERT INTO `techtrain_db`.`users` (`name`, `x_token`) VALUES (?, ?)", userReq.Name, userReq.XToken)
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

	result := s.Db.First(&s.User, "x_token = ?", token)
	c.JSON(http.StatusOK, gin.H{"name": result.Value.(*data.User).Name})
}

func (s *UserService) Update(c *gin.Context) {
	if err := c.BindJSON(&s.User); err != nil {
		log.Fatal(err)
	}
	token := c.GetHeader("x-token")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Token Required"})
		return
	}

	userReq := data.User{Name: s.User.Name, XToken: token}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/user/get", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("x-token", userReq.XToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&s.User)
	if err != nil {
		panic(err)
	}

	newName := userReq.Name
	oldName := s.User.Name

	s.Db.Exec("UPDATE `techtrain_db`.`users` SET `name` = ? WHERE (`name` = ?) and (`x_token` = ?)", newName, oldName, token)
	c.JSON(http.StatusOK, nil)
}
