package main

import (
	"chatapp/internal/config"
	"chatapp/internal/postgres"
	"chatapp/internal/router"
	"log"
	"net/http"
)

func main() {

	config.Load()
	postgres.Init()
	router := router.NewRouter()

	log.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", router)
}
