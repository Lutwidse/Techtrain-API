package techtrain_api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	Name   string `gorm:"column:name"`
	xToken string `gorm:"column:x_token"`
}

func (u *User) UserCreate(c *gin.Context) {
	db, err := gorm.Open("mysql", SqlArgs)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.LogMode(true)

	var user User
	if err := c.BindJSON(&user); err != nil {
		log.Fatal(err)
	}
	token := uuid.New().String()
	userReq := User{Name: user.Name, xToken: token}

	result := db.Exec("INSERT INTO `techtrain_db`.`users` (`name`, `x_token`) VALUES (?, ?)", userReq.Name, userReq.xToken)
	c.JSON(http.StatusOK, gin.H{"token": result.Value.(*User).xToken})
}

func (u *User) UserGet(c *gin.Context) {
	db, err := gorm.Open("mysql", SqlArgs)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.LogMode(true)
	var user User

	token := c.GetHeader("x-token")
	result := db.First(&user, "x_token = ?", token)
	c.JSON(http.StatusOK, gin.H{"name": result.Value.(*User).Name})
}

func (u *User) UserUpdate(c *gin.Context) {
	db, err := gorm.Open("mysql", SqlArgs)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.LogMode(true)

	var user User
	if err := c.BindJSON(&user); err != nil {
		log.Fatal(err)
	}
	token := c.GetHeader("x-token")
	userReq := User{Name: user.Name, xToken: token}

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

	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		panic(err)
	}

	newName := userReq.Name
	oldName := user.Name

	db.Exec("UPDATE `techtrain_db`.`users` SET `name` = ? WHERE (`name` = ?) and (`x_token` = ?)", newName, oldName, token)
	c.JSON(http.StatusOK, nil)
}
