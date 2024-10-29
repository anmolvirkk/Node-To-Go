const express = require('express');
const axios = require('axios');
require('dotenv').config();

const app = express();
const port = process.env.PORT || 3000;

app.use(express.json());

// Default coordinates for Brisbane
const BRISBANE_LAT = -27.4705;
const BRISBANE_LON = 153.0260;

// Middleware to validate coordinates
const validateCoordinates = (req, res, next) => {
    let { latitude, longitude } = req.query;
    
    // Use Brisbane coordinates if none provided
    if (!latitude && !longitude) {
        latitude = BRISBANE_LAT;
        longitude = BRISBANE_LON;
        req.query.latitude = latitude;
        req.query.longitude = longitude;
    } else if (!latitude || !longitude) {
        return res.status(400).json({ 
            error: 'Missing parameters',
            message: 'Both latitude and longitude must be provided if specifying location'
        });
    }

    const lat = parseFloat(latitude);
    const lon = parseFloat(longitude);

    if (isNaN(lat) || lat < -90 || lat > 90) {
        return res.status(400).json({ 
            error: 'Invalid latitude',
            message: 'Latitude must be between -90 and 90'
        });
    }

    if (isNaN(lon) || lon < -180 || lon > 180) {
        return res.status(400).json({ 
            error: 'Invalid longitude',
            message: 'Longitude must be between -180 and 180'
        });
    }

    next();
};

// Route to get current weather
app.get('/api/weather/current', validateCoordinates, async (req, res) => {
    try {
        const { latitude, longitude } = req.query;
        
        const response = await axios.get('https://api.open-meteo.com/v1/forecast', {
            params: {
                latitude,
                longitude,
                current: [
                    'temperature_2m',
                    'relative_humidity_2m',
                    'apparent_temperature',
                    'precipitation',
                    'rain',
                    'wind_speed_10m',
                    'wind_direction_10m',
                    'weather_code'
                ].join(','),
                timezone: 'auto'
            }
        });

        res.json({
            success: true,
            data: response.data.current,
            units: response.data.current_units
        });
    } catch (error) {
        console.error('Error fetching current weather:', error.message);
        res.status(500).json({ 
            error: 'Internal server error',
            message: 'Failed to fetch weather data'
        });
    }
});

// Route to get forecast
app.get('/api/weather/forecast', validateCoordinates, async (req, res) => {
    try {
        const { latitude, longitude } = req.query;
        const days = parseInt(req.query.days) || 7;
        
        if (days < 1 || days > 16) {
            return res.status(400).json({ 
                error: 'Invalid days parameter',
                message: 'Days must be between 1 and 16'
            });
        }

        const response = await axios.get('https://api.open-meteo.com/v1/forecast', {
            params: {
                latitude,
                longitude,
                daily: [
                    'temperature_2m_max',
                    'temperature_2m_min',
                    'precipitation_sum',
                    'rain_sum',
                    'precipitation_probability_max',
                    'wind_speed_10m_max',
                    'weather_code'
                ].join(','),
                timezone: 'auto',
                forecast_days: days
            }
        });

        res.json({
            success: true,
            data: response.data.daily,
            units: response.data.daily_units
        });
    } catch (error) {
        console.error('Error fetching forecast:', error.message);
        res.status(500).json({ 
            error: 'Internal server error',
            message: 'Failed to fetch forecast data'
        });
    }
});

// Error handling for invalid routes
app.use((req, res) => {
    res.status(404).json({ 
        error: 'Not Found',
        message: 'The requested resource does not exist'
    });
});

app.listen(port, () => {
    console.log(`Weather API server is running on port ${port}`);
}); 