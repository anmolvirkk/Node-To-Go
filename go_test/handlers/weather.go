package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCurrentWeather(c *gin.Context) {
	latitude := c.Query("latitude")
	longitude := c.Query("longitude")

	// Check if we have values set in middleware
	if lat, exists := c.Get("latitude"); exists {
		latitude = lat.(string)
	}
	if lon, exists := c.Get("longitude"); exists {
		longitude = lon.(string)
	}

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current=temperature_2m,relative_humidity_2m,apparent_temperature,precipitation,rain,wind_speed_10m,wind_direction_10m,weather_code&timezone=auto",
		latitude, longitude)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching current weather: %v", err)
		c.JSON(500, gin.H{
			"error":   "Internal server error",
			"message": "Failed to fetch weather data",
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		c.JSON(500, gin.H{
			"error":   "Internal server error",
			"message": "Failed to fetch weather data",
		})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error parsing response: %v", err)
		c.JSON(500, gin.H{
			"error":   "Internal server error",
			"message": "Failed to parse weather data",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    result["current"],
		"units":   result["current_units"],
	})
}

func GetForecast(c *gin.Context) {
	latitude := c.Query("latitude")
	longitude := c.Query("longitude")
	days := c.DefaultQuery("days", "7")

	// Check if we have values set in middleware
	if lat, exists := c.Get("latitude"); exists {
		latitude = lat.(string)
	}
	if lon, exists := c.Get("longitude"); exists {
		longitude = lon.(string)
	}

	daysInt, err := strconv.Atoi(days)
	if err != nil || daysInt < 1 || daysInt > 16 {
		c.JSON(400, gin.H{
			"error":   "Invalid days parameter",
			"message": "Days must be between 1 and 16",
		})
		return
	}

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&daily=temperature_2m_max,temperature_2m_min,precipitation_sum,rain_sum,precipitation_probability_max,wind_speed_10m_max,weather_code&timezone=auto&forecast_days=%d",
		latitude, longitude, daysInt)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching forecast: %v", err)
		c.JSON(500, gin.H{
			"error":   "Internal server error",
			"message": "Failed to fetch forecast data",
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		c.JSON(500, gin.H{
			"error":   "Internal server error",
			"message": "Failed to fetch forecast data",
		})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error parsing response: %v", err)
		c.JSON(500, gin.H{
			"error":   "Internal server error",
			"message": "Failed to parse forecast data",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    result["daily"],
		"units":   result["daily_units"],
	})
} 