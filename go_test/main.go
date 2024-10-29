package main

import (
	"go_test/handlers"
	"go_test/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Set up router
	r := gin.Default()

	// Routes
	r.GET("/api/weather/current", middleware.ValidateCoordinates(), handlers.GetCurrentWeather)
	r.GET("/api/weather/forecast", middleware.ValidateCoordinates(), handlers.GetForecast)

	// Handle 404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "Not Found",
			"message": "The requested resource does not exist",
		})
	})

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Start server
	log.Printf("Weather API server is running on port %s", port)
	r.Run(":" + port)
} 