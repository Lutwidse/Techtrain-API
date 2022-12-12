package techtrain_api

import (
	"github.com/gin-gonic/gin"
)

type User struct {
	Name    string
	xToken string `gorm:"column:x_token"`
}

func (u *User) UserCreate(c *gin.Context) {

}

func (u *User) UserGet(c *gin.Context) {

}

func (u *User) UserUpdate(c *gin.Context) {

}
