package main

import (
	"html/template"
	"log"

	"github.com/NenfuAT/24AuthorizationServer/model"
	"github.com/NenfuAT/24AuthorizationServer/router"
)

func main() {

	var err error
	model.Templates["login"], err = template.ParseFiles("front/login.html")
	if err != nil {
		log.Fatal(err)
	}
	router.Init()
}
