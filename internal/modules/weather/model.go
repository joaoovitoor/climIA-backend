package weather

import (
	"time"
)

type Weather struct {
	ID                uint      `json:"id"`
	Data              time.Time `json:"data"`
	TemperaturaMinima float64   `json:"temperatura_minima"`
	TemperaturaMaxima float64   `json:"temperatura_maxima"`
	TemperaturaMedia  float64   `json:"temperatura_media"`
	Cidade            string    `json:"cidade"`
	Estado            string    `json:"estado"`
	Precipitacao      float64   `json:"precipitacao"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type WeatherRequest struct {
	Cidade     string `json:"cidade" query:"cidade"`
	Estado     string `json:"estado" query:"estado"`
	Data       string `json:"data" query:"data"`
	DataInicio string `json:"datainicio" query:"datainicio"`
	DataFim    string `json:"datafim" query:"datafim"`
}

type WeatherResponse struct {
	Cidade            string  `json:"cidade"`
	Estado            string  `json:"estado"`
	Data              string  `json:"data"`
	TemperaturaMinima float64 `json:"temperatura_minima"`
	TemperaturaMaxima float64 `json:"temperatura_maxima"`
	TemperaturaMedia  float64 `json:"temperatura_media"`
	Precipitacao      float64 `json:"precipitacao"`
} 