package model

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
