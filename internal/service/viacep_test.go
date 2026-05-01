package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/temperaturaPorCep/internal/model"
)

func TestViaCEP_GetLocation_Success(t *testing.T) {
	// Criar servidor HTTP mock
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/01001000/json/"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		resp := model.ViaCEPResponse{
			Cep:         "01001-000",
			Logradouro:  "Praça da Sé",
			Bairro:      "Sé",
			Localidade:  "São Paulo",
			Uf:          "SP",
			Ibge:        "3550308",
			Gia:         "9",
			Ddd:         "11",
			Siafi:       "7107",
			Erro:       false,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	s := &viacepService{
		client:  &http.Client{},
		baseURL: server.URL,
	}

	location, err := s.GetLocation("01001000")
	if err != nil {
		t.Fatalf("GetLocation returned error: %v", err)
	}

	if location.Localidade != "São Paulo" {
		t.Errorf("Expected city 'São Paulo', got '%s'", location.Localidade)
	}
	if location.Uf != "SP" {
		t.Errorf("Expected state 'SP', got '%s'", location.Uf)
	}
	if location.Cep != "01001-000" {
		t.Errorf("Expected CEP '01001-000', got '%s'", location.Cep)
	}
}

func TestViaCEP_GetLocation_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := model.ViaCEPResponse{Erro: true}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	s := &viacepService{
		client:  &http.Client{},
		baseURL: server.URL,
	}

	_, err := s.GetLocation("00000000")
	if err == nil {
		t.Error("Expected error for not found CEP, got nil")
	}
}

func TestViaCEP_GetLocation_InvalidServerResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	s := &viacepService{
		client:  &http.Client{},
		baseURL: server.URL,
	}

	_, err := s.GetLocation("99999999")
	if err == nil {
		t.Error("Expected error for 500 response, got nil")
	}
}
