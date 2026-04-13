package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pickkinsley/project2/backend/db"
	"github.com/pickkinsley/project2/backend/handlers"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "packsmart_user:Dogs1234@tcp(localhost)/packsmart?parseTime=true"
	}

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to MySQL")

	queries := db.New(sqlDB)
	h := handlers.NewHandler(sqlDB, queries)

	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/health", handlers.Health)
		api.POST("/trips", h.CreateTrip)
		api.GET("/trips/:uuid", h.GetTrip)
		api.PATCH("/trips/:uuid/items/:itemId", h.UpdateItemCheckbox)
	}

	log.Println("PackSmart backend starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
