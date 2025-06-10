package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type GeoData struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func GetCoordinates(ctx context.Context, city string, state string, apiKey string) (GeoData, error) {
	var result []GeoData
	url := fmt.Sprintf("https://api.openweathermap.org/geo/1.0/direct?q=%s,%s,BR&limit=1&appid=%s", city, state, apiKey)
	log.Printf("Geocoding API URL: %s", url)
	
	body, err := FetchAPI(ctx, url)
	if err != nil {
		log.Printf("Error fetching coordinates: %v", err)
		return GeoData{}, err
	}
	
	log.Printf("Geocoding API Response: %s", body.body)
	
	if err := json.Unmarshal([]byte(body.body), &result); err != nil {
		log.Printf("Error unmarshaling coordinates: %v", err)
		return GeoData{}, err
	}
	
	if len(result) == 0 {
		return GeoData{}, fmt.Errorf("city not found")
	}
	
	return result[0], nil
}

func GetWeatherAPI(ctx context.Context, lat, lon float64, apiKey string) (WeatherData, error) {
	var result WeatherData
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%.6f&lon=%.6f&appid=%s", lat, lon, apiKey)
	log.Printf("Weather API URL: %s", url)
	
	body, err := FetchAPI(ctx, url)
	if err != nil {
		log.Printf("Error fetching weather data: %v", err)
		return WeatherData{}, err
	}
	
	log.Printf("Weather API Response: %s", body.body)
	
	if err := json.Unmarshal([]byte(body.body), &result); err != nil {
		log.Printf("Error unmarshaling weather data: %v", err)
		return WeatherData{}, err
	}
	return result, nil
}

func ParseTemperatureResponse(weather WeatherData) Temperature {
	log.Printf("Raw weather data: %+v", weather)
	
	kelvin := weather.Main.Temp
	log.Printf("Temperature in Kelvin: %f", kelvin)
	
	celsius := kelvin - 273.15
	log.Printf("Temperature in Celsius: %f", celsius)
	
	fahrenheit := celsius*1.8 + 32
	log.Printf("Temperature in Fahrenheit: %f", fahrenheit)
	
	temp := Temperature{
		Celsius:    celsius,
		Fahrenheit: fahrenheit,
		Kelvin:     kelvin,
	}
	
	log.Printf("Final temperature object: %+v", temp)
	return temp
} 