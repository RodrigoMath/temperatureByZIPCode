package converter

import "math"

// CelsiusToFahrenheit converte Celsius para Fahrenheit
func CelsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32
}

// CelsiusToKelvin converte Celsius para Kelvin
func CelsiusToKelvin(c float64) float64 {
	return c + 273.15
}

// ConvertAll converte Celsius para Fahrenheit e Kelvin
func ConvertAll(celsius float64) (c, f, k float64) {
	f = CelsiusToFahrenheit(celsius)
	k = CelsiusToKelvin(celsius)
	return celsius, f, k
}

// RoundTo1Decimal arredonda para 1 casa decimal
func RoundTo1Decimal(value float64) float64 {
	return math.Round(value*10) / 10
}
