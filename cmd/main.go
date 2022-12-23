package main

import (
	"io/ioutil"
	"log"

	api "github.com/Lutwidse/Techtrain-API"
	yaml "gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SQL struct {
	User string `yaml:"user"`
	Pw   string `yaml:"pw"`
	Db   string `yaml:"db"`
}

func main() {
	file, err := ioutil.ReadFile("../cmd/config/sql.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var sql SQL

	err = yaml.Unmarshal(file, &sql)
	if err != nil {
		log.Fatal(err)
	}

	loginInfo := (sql.User + ":" + sql.Pw + "/" + sql.Db)
	db, err := gorm.Open(mysql.Open(loginInfo), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	srv := api.NewTechtrainServer(db)
	srv.Start()
}
