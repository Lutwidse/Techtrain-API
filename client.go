package api

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type TechtrainClient struct {
	user      *User
	gacha     *Gacha
	character *Character
}

func NewTechtrainClient() *TechtrainClient {
	return &TechtrainClient{user: &User{}, gacha: &Gacha{}, character: &Character{}}
}

func (c *TechtrainClient) Test() {
	db, err := sql.Open("mysql", "root:$A!c6i7xC!93%@FZ@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("Success!")
}
