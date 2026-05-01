package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/temperaturaPorCep/internal/model"
	"github.com/temperaturaPorCep/internal/service"
	"github.com/temperaturaPorCep/pkg/converter"
)

// WeatherHandler lida com as requisições de temperatura por CEP
type WeatherHandler struct {
	viacepService  service.ViaCEPService
	weatherService service.WeatherService
}

// NewWeatherHandler cria um novo handler
func NewWeatherHandler(viacepService service.ViaCEPService, weatherService service.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		viacepService:  viacepService,
		weatherService: weatherService,
	}
}

// GetTemperatureByCEP retorna a temperatura para um CEP
func (h *WeatherHandler) GetTemperatureByCEP(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	// Validação do CEP: deve ter exatamente 8 dígitos
	if len(cep) != 8 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity) // 422
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "invalid zipcode"})
		return
	}

	// Busca a localização pelo CEP
	location, err := h.viacepService.GetLocation(cep)
	if err != nil {
		if errors.Is(err, service.ErrZipCodeNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound) // 404
			json.NewEncoder(w).Encode(model.ErrorResponse{Message: "can not find zipcode"})
			return
		}
		// Erro interno
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "internal server error"})
		return
	}

	// Busca a temperatura pela cidade
	city := location.Localidade
	tempC, err := h.weatherService.GetTemperature(city)
	if err != nil {
		// Qualquer erro do weather service retorna 500
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Message: "internal server error"})
		return
	}

	// Converte temperaturas
	_, tempF, tempK := converter.ConvertAll(tempC)

	// Formata resposta com 1 casa decimal
	response := model.TemperatureResponse{
		TempC: converter.RoundTo1Decimal(tempC),
		TempF: converter.RoundTo1Decimal(tempF),
		TempK: converter.RoundTo1Decimal(tempK),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	json.NewEncoder(w).Encode(response)
}
