package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/temperaturaPorCep/internal/model"
	"github.com/temperaturaPorCep/internal/service"
)

// Mock ViaCEP
type mockViaCEP struct {
	location *model.ViaCEPResponse
	err      error
}

func (m *mockViaCEP) GetLocation(cep string) (*model.ViaCEPResponse, error) {
	return m.location, m.err
}

// Mock Weather
type mockWeather struct {
	tempC float64
	err    error
}

func (m *mockWeather) GetTemperature(city string) (float64, error) {
	return m.tempC, m.err
}

func TestWeatherHandler_GetTemperatureByCEP_Success(t *testing.T) {
	mockViacep := &mockViaCEP{
		location: &model.ViaCEPResponse{
			Localidade: "São Paulo",
		},
	}
	mockWeather := &mockWeather{
		tempC: 28.5,
	}

	h := NewWeatherHandler(mockViacep, mockWeather)

	req := httptest.NewRequest("GET", "/01001000", nil)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/{cep}", h.GetTemperatureByCEP)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp model.TemperatureResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.TempC != 28.5 {
		t.Errorf("Expected TempC 28.5, got %v", resp.TempC)
	}
	if resp.TempF != 83.3 {
		t.Errorf("Expected TempF 83.3, got %v", resp.TempF)
	}
	if resp.TempK != 301.7 {
		t.Errorf("Expected TempK 301.7, got %v", resp.TempK)
	}
}

func TestWeatherHandler_GetTemperatureByCEP_InvalidFormat_Short(t *testing.T) {
	h := NewWeatherHandler(&mockViaCEP{}, &mockWeather{})

	req := httptest.NewRequest("GET", "/123", nil)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/{cep}", h.GetTemperatureByCEP)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d", w.Code)
	}

	var errResp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}
	if errResp.Message != "invalid zipcode" {
		t.Errorf("Expected message 'invalid zipcode', got '%s'", errResp.Message)
	}
}

func TestWeatherHandler_GetTemperatureByCEP_InvalidFormat_Long(t *testing.T) {
	h := NewWeatherHandler(&mockViaCEP{}, &mockWeather{})

	req := httptest.NewRequest("GET", "/123456789", nil)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/{cep}", h.GetTemperatureByCEP)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d", w.Code)
	}
}

func TestWeatherHandler_GetTemperatureByCEP_NotFound(t *testing.T) {
	mockViacep := &mockViaCEP{
		err: service.ErrZipCodeNotFound,
	}
	h := NewWeatherHandler(mockViacep, &mockWeather{})

	req := httptest.NewRequest("GET", "/00000000", nil)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/{cep}", h.GetTemperatureByCEP)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	var errResp model.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}
	if errResp.Message != "can not find zipcode" {
		t.Errorf("Expected message 'can not find zipcode', got '%s'", errResp.Message)
	}
}

func TestWeatherHandler_GetTemperatureByCEP_WeatherError(t *testing.T) {
	mockViacep := &mockViaCEP{
		location: &model.ViaCEPResponse{
			Localidade: "São Paulo",
		},
	}
	mockWeather := &mockWeather{
		err: errors.New("weather service error"),
	}
	h := NewWeatherHandler(mockViacep, mockWeather)

	req := httptest.NewRequest("GET", "/01001000", nil)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/{cep}", h.GetTemperatureByCEP)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}
