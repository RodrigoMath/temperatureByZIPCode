package model

import (
	"encoding/json"
	"testing"
)

func TestViaCEPResponse_UnmarshalJSON_WithBool(t *testing.T) {
	data := []byte(`{"cep":"01001-000","localidade":"São Paulo","erro":false}`)
	var v ViaCEPResponse
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if v.Erro != false {
		t.Errorf("Expected Erro=false, got %v", v.Erro)
	}
	if v.Localidade != "São Paulo" {
		t.Errorf("Expected Localidade='São Paulo', got '%s'", v.Localidade)
	}
}

func TestViaCEPResponse_UnmarshalJSON_WithTrueBool(t *testing.T) {
	data := []byte(`{"erro":true}`)
	var v ViaCEPResponse
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if v.Erro != true {
		t.Errorf("Expected Erro=true, got %v", v.Erro)
	}
}

func TestViaCEPResponse_UnmarshalJSON_WithString(t *testing.T) {
	// ViaCEP às vezes retorna "true" como string
	data := []byte(`{"erro":"true"}`)
	var v ViaCEPResponse
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if v.Erro != true {
		t.Errorf("Expected Erro=true, got %v", v.Erro)
	}
}

func TestViaCEPResponse_UnmarshalJSON_WithFalseString(t *testing.T) {
	data := []byte(`{"erro":"false"}`)
	var v ViaCEPResponse
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if v.Erro != false {
		t.Errorf("Expected Erro=false, got %v", v.Erro)
	}
}

func TestViaCEPResponse_UnmarshalJSON_WithNumber(t *testing.T) {
	data := []byte(`{"erro":1}`)
	var v ViaCEPResponse
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if v.Erro != true {
		t.Errorf("Expected Erro=true, got %v", v.Erro)
	}
}
