package main

import (
	"html/template"
	"log"

	"github.com/kajiLabTeam/mr-platform-authorization-server/model"
	"github.com/kajiLabTeam/mr-platform-authorization-server/router"
)

func main() {

	var err error
	model.Templates["login"], err = template.ParseFiles("front/login.html")
	if err != nil {
		log.Fatal(err)
	}
	router.Init()
}
