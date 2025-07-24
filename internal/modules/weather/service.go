package weather

import (
	"fmt"
	"time"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CalculateForecast(req WeatherRequest) ([]WeatherResponse, error) {
	if req.Cidade == "" || req.Estado == "" {
		return nil, fmt.Errorf("cidade e estado são obrigatórios")
	}

	var dataInicio, dataFim *time.Time

	if req.Data != "" {
		data, err := time.Parse("2006-01-02", req.Data)
		if err != nil {
			return nil, fmt.Errorf("formato de data inválido. Use YYYY-MM-DD")
		}
		weather, err := s.repository.GetWeatherByDate(req.Cidade, req.Estado, data)
		if err != nil {
			return nil, fmt.Errorf("dados não encontrados para a data especificada")
		}
		return []WeatherResponse{s.convertToResponse(*weather)}, nil
	}

	if req.DataInicio != "" {
		data, err := time.Parse("2006-01-02", req.DataInicio)
		if err != nil {
			return nil, fmt.Errorf("formato de data de início inválido. Use YYYY-MM-DD")
		}
		dataInicio = &data
	}

	if req.DataFim != "" {
		data, err := time.Parse("2006-01-02", req.DataFim)
		if err != nil {
			return nil, fmt.Errorf("formato de data de fim inválido. Use YYYY-MM-DD")
		}
		dataFim = &data
	}

	weatherData, err := s.repository.GetWeatherData(req.Cidade, req.Estado, dataInicio, dataFim)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar dados meteorológicos: %v", err)
	}

	if len(weatherData) == 0 {
		return nil, fmt.Errorf("nenhum dado encontrado para %s, %s", req.Cidade, req.Estado)
	}

	var responses []WeatherResponse
	for _, weather := range weatherData {
		responses = append(responses, s.convertToResponse(weather))
	}

	return responses, nil
}

func (s *Service) convertToResponse(weather Weather) WeatherResponse {
	return WeatherResponse{
		Cidade:            weather.Cidade,
		Estado:            weather.Estado,
		Data:              weather.Data.Format("2006-01-02"),
		TemperaturaMinima: weather.TemperaturaMinima,
		TemperaturaMaxima: weather.TemperaturaMaxima,
		TemperaturaMedia:  weather.TemperaturaMedia,
		Precipitacao:      weather.Precipitacao,
	}
} 