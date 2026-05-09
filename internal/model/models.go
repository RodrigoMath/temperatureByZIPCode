package model

import "encoding/json"

// ViaCEPResponse representa a resposta da API ViaCEP
type ViaCEPResponse struct {
	Cep          string `json:"cep"`
	Logradouro   string `json:"logradouro,omitempty"`
	Complemento  string `json:"complemento,omitempty"`
	Bairro       string `json:"bairro,omitempty"`
	Localidade   string `json:"localidade"` // Cidade
	Uf           string `json:"uf,omitempty"`
	Ibge         string `json:"ibge,omitempty"`
	Gia          string `json:"gia,omitempty"`
	Ddd          string `json:"ddd,omitempty"`
	Siafi        string `json:"siafi,omitempty"`
	Erro        bool   `json:"erro"` // ViaCEP retorna true para CEPs não encontrados
}

// UnmarshalJSON customizado para lidar com 'erro' como string ou bool
func (v *ViaCEPResponse) UnmarshalJSON(data []byte) error {
	// Struct auxiliar para evitar recursão
	type Alias ViaCEPResponse
	aux := &struct {
		Erro interface{} `json:"erro"`
		*Alias
	}{
		Alias: (*Alias)(v),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Converter Erro para bool
	switch e := aux.Erro.(type) {
	case bool:
		v.Erro = e
	case string:
		// ViaCEP pode retornar "true" ou "false" como string
		v.Erro = e == "true" || e == "1"
	case float64:
		// JSON números também podem aparecer
		v.Erro = e == 1
	case nil:
		v.Erro = false
	default:
		// Se for outro tipo, assumir false
		v.Erro = false
	}

	return nil
}

// TemperatureResponse representa a resposta formatada de temperatura
type TemperatureResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// ErrorResponse representa mensagens de erro da API
type ErrorResponse struct {
	Message string `json:"error"`
}
