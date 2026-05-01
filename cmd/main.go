package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/temperaturaPorCep/internal/handler"
	"github.com/temperaturaPorCep/internal/service"
)

func main() {
	r := chi.NewRouter()

	// Inicializa serviços
	viacepService := service.NewViaCEPService()

	// Open-Meteo API - serviço de clima gratuito, sem necessidade de chave API
	weatherService := service.NewWeatherService("")

	// Inicializa handler
	weatherHandler := handler.NewWeatherHandler(viacepService, weatherService)

	// Rotas
	r.Get("/{cep}", weatherHandler.GetTemperatureByCEP)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
