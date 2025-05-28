package main

import (
	"chatapp/internal/config"
	"chatapp/internal/postgres"
	"chatapp/internal/router"
	"fmt"
	"log"
	"net/http"
)

func main() {

	config.Load()
	postgres.Init()
	router := router.NewRouter()

	log.Printf("Listening on port %v...", config.App.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", config.App.Port), router)
}
