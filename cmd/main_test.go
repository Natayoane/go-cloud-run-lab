package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ntayoane/go-cloud-run-lab/configs"
	"github.com/ntayoane/go-cloud-run-lab/internal"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("OPENWEATHERMAP_API_KEY", "test-api-key")
	code := m.Run()
	os.Unsetenv("OPENWEATHERMAP_API_KEY")
	os.Exit(code)
}

func TestValidateCEP(t *testing.T) {
	tests := []struct {
		name    string
		cep     string
		wantErr bool
	}{
		{
			name:    "valid CEP",
			cep:     "89010-904",
			wantErr: false,
		},
		{
			name:    "invalid CEP - wrong length",
			cep:     "89010",
			wantErr: true,
		},
		{
			name:    "invalid CEP - non-numeric",
			cep:     "89010-abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := internal.ValidateCEP(tt.cep)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTemperatureConversion(t *testing.T) {
	// OpenWeatherMap API returns temperature in Kelvin, so we use 298.15K (25°C)
	weatherData := internal.WeatherData{
		Main: struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		}{
			Temp:      298.15, // 25°C in Kelvin
			FeelsLike: 297.15, // 24°C in Kelvin
			TempMin:   293.15, // 20°C in Kelvin
			TempMax:   303.15, // 30°C in Kelvin
			Pressure:  1013,
			Humidity:  65,
		},
	}

	temp := internal.ParseTemperatureResponse(weatherData)

	assert.Equal(t, 25.0, temp.Celsius)
	assert.Equal(t, 77.0, temp.Fahrenheit)
	assert.Equal(t, 298.15, temp.Kelvin)
}

func TestHandleTemperatureRequest(t *testing.T) {
	config := configs.Config{
		OpenWeatherMapAPIKey: "test-api-key",
	}

	server := httptest.NewServer(http.HandlerFunc(internal.HandleTemperatureRequest(config)))
	defer server.Close()

	tests := []struct {
		name           string
		cep            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "missing CEP",
			cep:            "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "CEP is required",
		},
		{
			name:           "invalid CEP format",
			cep:            "123",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   "invalid zipcode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := server.URL + "/temperature"
			if tt.cep != "" {
				url += "?cep=" + tt.cep
			}

			resp, err := http.Get(url)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), tt.expectedBody)
		})
	}
}
