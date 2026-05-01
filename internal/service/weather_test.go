package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherService_GetTemperature_Success(t *testing.T) {
	geoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"results": []map[string]interface{}{
				{
					"name":      "São Paulo",
					"latitude":  -23.55,
					"longitude": -46.63,
					"country":   "Brazil",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer geoServer.Close()

	weatherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"current": map[string]interface{}{
				"temperature_2m": 28.5,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer weatherServer.Close()

	s := &weatherService{
		client:  &http.Client{},
		baseURL: weatherServer.URL,
		geoURL:  geoServer.URL,
	}

	temp, err := s.GetTemperature("São Paulo")
	if err != nil {
		t.Fatalf("GetTemperature returned error: %v", err)
	}

	if temp != 28.5 {
		t.Errorf("Expected temperature 28.5, got %v", temp)
	}
}

func TestWeatherService_GetTemperature_CityNotFound(t *testing.T) {
	geoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"results": []map[string]interface{}{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer geoServer.Close()

	s := &weatherService{
		client:  &http.Client{},
		geoURL:  geoServer.URL,
		baseURL: "http://example.com",
	}

	_, err := s.GetTemperature("CidadeInexistente12345")
	if err == nil {
		t.Error("Expected error for city not found, got nil")
	}
}

func TestWeatherService_GetTemperature_GeoAPIFailure(t *testing.T) {
	geoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer geoServer.Close()

	s := &weatherService{
		client:  &http.Client{},
		geoURL:  geoServer.URL,
		baseURL: "http://example.com",
	}

	_, err := s.GetTemperature("São Paulo")
	if err == nil {
		t.Error("Expected error for geocoding API failure, got nil")
	}
}
