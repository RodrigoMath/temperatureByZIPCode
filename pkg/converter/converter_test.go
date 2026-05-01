package converter

import (
	"math"
	"testing"
)

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		c    float64
		want float64
	}{
		{0, 32},
		{100, 212},
		{28.5, 83.3},
		{-40, -40},
		{37, 98.6},
	}
	for _, tt := range tests {
		got := CelsiusToFahrenheit(tt.c)
		if math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("CelsiusToFahrenheit(%v) = %v, want %v", tt.c, got, tt.want)
		}
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		c    float64
		want float64
	}{
		{0, 273.15},
		{100, 373.15},
		{28.5, 301.65},
		{-273.15, 0},
	}
	for _, tt := range tests {
		got := CelsiusToKelvin(tt.c)
		if got != tt.want {
			t.Errorf("CelsiusToKelvin(%v) = %v, want %v", tt.c, got, tt.want)
		}
	}
}

func TestRoundTo1Decimal(t *testing.T) {
	tests := []struct {
		value    float64
		expected float64
	}{
		{28.56, 28.6},
		{28.54, 28.5},
		{83.30, 83.3},
		{301.65, 301.7},
		{0.05, 0.1},
		{0.04, 0.0},
	}
	for _, tt := range tests {
		got := RoundTo1Decimal(tt.value)
		if got != tt.expected {
			t.Errorf("RoundTo1Decimal(%v) = %v, want %v", tt.value, got, tt.expected)
		}
	}
}
