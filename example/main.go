package main

import (
	api "github.com/Lutwidse/Techtrain-API"
)

func main() {
	c := api.NewTechtrainClient()
	c.Server()
}