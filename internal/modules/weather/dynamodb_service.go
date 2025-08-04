package weather

import (
	"fmt"
	"time"
)

type DynamoDBService struct {
	repository *DynamoDBRepository
}

func NewDynamoDBService(repository *DynamoDBRepository) *DynamoDBService {
	return &DynamoDBService{repository: repository}
}

func (s *DynamoDBService) GetProcessedForecast(req WeatherRequest) ([]WeatherResponse, error) {
	if req.Cidade == "" || req.Estado == "" {
		return nil, fmt.Errorf("cidade e estado são obrigatórios")
	}

	if req.Data != "" {
		data, err := time.Parse("2006-01-02", req.Data)
		if err != nil {
			return nil, fmt.Errorf("formato de data inválido. Use YYYY-MM-DD")
		}

		forecast, err := s.repository.GetWeatherByDate(req.Cidade, req.Estado, data)
		if err != nil {
			return nil, err
		}

		return []WeatherResponse{*forecast}, nil
	}

	if req.DataInicio != "" && req.DataFim != "" {
		dataInicio, err := time.Parse("2006-01-02", req.DataInicio)
		if err != nil {
			return nil, fmt.Errorf("formato de data início inválido. Use YYYY-MM-DD")
		}

		dataFim, err := time.Parse("2006-01-02", req.DataFim)
		if err != nil {
			return nil, fmt.Errorf("formato de data fim inválido. Use YYYY-MM-DD")
		}

		return s.repository.GetWeatherDataForPeriod(req.Cidade, req.Estado, dataInicio, dataFim)
	}

	return nil, fmt.Errorf("deve fornecer data ou intervalo de datas")
}
