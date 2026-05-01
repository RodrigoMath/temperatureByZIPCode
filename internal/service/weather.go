package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Erros específicos dos serviços
var (
	ErrZipCodeNotFound   = errors.New("zipcode not found")
	ErrCityNotFound      = errors.New("city not found in weather API")
	ErrInvalidAPIKey     = errors.New("invalid WeatherAPI key")
	ErrRateLimitExceeded = errors.New("weather API rate limit exceeded")
)

// WeatherService interface para buscar temperatura
type WeatherService interface {
	GetTemperature(city string) (float64, error)
}

// openMeteoResponse representa a resposta da API Open-Meteo para previsão
type openMeteoResponse struct {
	Current struct {
		Temperature2m float64 `json:"temperature_2m"`
	} `json:"current"`
}

// openMeteoGeoResponse representa a resposta da API Open-Meteo para geocoding
type openMeteoGeoResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Country   string  `json:"country"`
	} `json:"results"`
}

type weatherService struct {
	client  *http.Client
	baseURL string
	geoURL  string
}

// NewWeatherService creates a new WeatherService instance
// Usa Open-Meteo API (gratuita, sem necessidade de chave API)
func NewWeatherService(_ string) WeatherService {
	return &weatherService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://api.open-meteo.com/v1/forecast",
		geoURL:  "https://geocoding-api.open-meteo.com/v1/search",
	}
}

func (s *weatherService) GetTemperature(city string) (float64, error) {
	// Primeiro, obter coordenadas via geocoding API
	geoURL := fmt.Sprintf("%s?name=%s&count=1", s.geoURL, url.QueryEscape(city))
	geoReq, err := http.NewRequest("GET", geoURL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create geocoding request: %w", err)
	}

	geoResp, err := s.client.Do(geoReq)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch geocoding: %w", err)
	}
	defer geoResp.Body.Close()

	if geoResp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("geocoding API returned status %d", geoResp.StatusCode)
	}

	var geoResult openMeteoGeoResponse
	if err := json.NewDecoder(geoResp.Body).Decode(&geoResult); err != nil {
		return 0, fmt.Errorf("failed to decode geocoding response: %w", err)
	}

	if len(geoResult.Results) == 0 {
		return 0, ErrCityNotFound
	}

	lat := geoResult.Results[0].Latitude
	lon := geoResult.Results[0].Longitude

	// Segundo, obter temperatura via previsão
	u, err := url.Parse(s.baseURL)
	if err != nil {
		return 0, err
	}

	query := u.Query()
	query.Set("latitude", fmt.Sprintf("%f", lat))
	query.Set("longitude", fmt.Sprintf("%f", lon))
	query.Set("current", "temperature_2m")
	query.Set("timezone", "auto")
	u.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var result openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode weather response: %w", err)
	}

	return result.Current.Temperature2m, nil
}
