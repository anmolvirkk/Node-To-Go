package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	BRISBANE_LAT = -27.4705
	BRISBANE_LON = 153.0260
)

func ValidateCoordinates() gin.HandlerFunc {
	return func(c *gin.Context) {
		latitude := c.Query("latitude")
		longitude := c.Query("longitude")

		// Use Brisbane coordinates if none provided
		if latitude == "" && longitude == "" {
			c.Set("latitude", strconv.FormatFloat(BRISBANE_LAT, 'f', -1, 64))
			c.Set("longitude", strconv.FormatFloat(BRISBANE_LON, 'f', -1, 64))
			latitude = strconv.FormatFloat(BRISBANE_LAT, 'f', -1, 64)
			longitude = strconv.FormatFloat(BRISBANE_LON, 'f', -1, 64)
		} else if latitude == "" || longitude == "" {
			c.JSON(400, gin.H{
				"error":   "Missing parameters",
				"message": "Both latitude and longitude must be provided if specifying location",
			})
			c.Abort()
			return
		}

		lat, err := strconv.ParseFloat(latitude, 64)
		if err != nil || lat < -90 || lat > 90 {
			c.JSON(400, gin.H{
				"error":   "Invalid latitude",
				"message": "Latitude must be between -90 and 90",
			})
			c.Abort()
			return
		}

		lon, err := strconv.ParseFloat(longitude, 64)
		if err != nil || lon < -180 || lon > 180 {
			c.JSON(400, gin.H{
				"error":   "Invalid longitude",
				"message": "Longitude must be between -180 and 180",
			})
			c.Abort()
			return
		}

		c.Next()
	}
} 