package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pickkinsley/project2/backend/handlers"
)

func main() {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/health", handlers.Health)
		api.POST("/trips", handlers.CreateTrip)
		api.GET("/trips/:uuid", handlers.GetTrip)
	}

	log.Println("PackSmart backend starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
