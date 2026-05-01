package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/temperaturaPorCep/internal/model"
)

// ViaCEPService interface para buscar localização por CEP
type ViaCEPService interface {
	GetLocation(cep string) (*model.ViaCEPResponse, error)
}

type viacepService struct {
	client  *http.Client
	baseURL string
}

// NewViaCEPService cria uma nova instância do serviço ViaCEP
func NewViaCEPService() ViaCEPService {
	return &viacepService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://viacep.com.br/ws",
	}
}

func (s *viacepService) GetLocation(cep string) (*model.ViaCEPResponse, error) {
	url := fmt.Sprintf("%s/%s/json/", s.baseURL, cep)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ViaCEP returned status %d", resp.StatusCode)
	}

	var result model.ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode ViaCEP response: %w", err)
	}

	// ViaCEP retorna {"erro": true} para CEPs não encontrados
	if result.Erro {
		return nil, ErrZipCodeNotFound
	}

	return &result, nil
}
