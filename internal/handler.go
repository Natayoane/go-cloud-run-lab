package internal

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ntayoane/go-cloud-run-lab/configs"
)

func HandleTemperatureRequest(config configs.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		cep := r.URL.Query().Get("cep")
		if cep == "" {
			http.Error(w, "CEP is required", http.StatusBadRequest)
			return
		}
		if err := ValidateCEP(cep); err != nil {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		
		viaCepData, err := GetViaCepAPI(ctx, cep)
		if err != nil {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}
		
		geoData, err := GetCoordinates(ctx, viaCepData.Localidade, viaCepData.Uf, config.OpenWeatherMapAPIKey)
		if err != nil {
			log.Printf("Error getting coordinates: %v", err)
			http.Error(w, "Error getting coordinates", http.StatusInternalServerError)
			return
		}
		
		log.Printf("Using coordinates - Latitude: %f, Longitude: %f", geoData.Lat, geoData.Lon)
		
		weatherData, err := GetWeatherAPI(ctx, geoData.Lat, geoData.Lon, config.OpenWeatherMapAPIKey)
		if err != nil {
			log.Printf("Error in GetWeatherAPI: %v", err)
			http.Error(w, "Error fetching weather data", http.StatusInternalServerError)
			return
		}
		temperature := ParseTemperatureResponse(weatherData)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(temperature)
	}
} 