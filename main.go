package main

import (
	"log"

	"url-shortener/routes"
)

func main() {
	router := routes.SetupRouter()

	log.Println("Server starting on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
