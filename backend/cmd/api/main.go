package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mvcris/maya-guessr/backend/internal/infrastructure/gorm"
	"github.com/mvcris/maya-guessr/backend/internal/interfaces/http/handlers"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatalf("DATABASE_URL is not set")
	}
	db, err := gorm.NewConnection(databaseUrl)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	router := gin.Default()

	// users routes
	userHandler := handlers.NewUserHandler(db, router)
	userHandler.SetupRoutes()

	// auth routes
	authHandler := handlers.NewAuthHandler(db, router)
	authHandler.SetupRoutes()

	// maps routes
	mapHandler := handlers.NewMapHandler(db, router)
	mapHandler.SetupRoutes()

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}