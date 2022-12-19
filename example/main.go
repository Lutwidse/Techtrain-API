package main

import (
	"log"
	api "github.com/Lutwidse/Techtrain-API"
	"github.com/jinzhu/gorm"
)

func main() {
	db, err := gorm.Open("mysql", "root:$A!c6i7xC!93%@FZ@/techtrain_db")
	if err != nil {
		log.Fatal(err)
	}

	c := api.NewTechtrainServer(db)
	c.Server()
}
