package weather

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetWeatherData(cidade, estado string, dataInicio, dataFim *time.Time) ([]Weather, error) {
	var weatherData []Weather
	query := r.db.Table("previsao_tempo").Where("cidade = ? AND estado = ?", cidade, estado)

	if dataInicio != nil && dataFim != nil {
		query = query.Where("data BETWEEN ? AND ?", dataInicio, dataFim)
	} else if dataInicio != nil {
		query = query.Where("data >= ?", dataInicio)
	} else if dataFim != nil {
		query = query.Where("data <= ?", dataFim)
	}

	err := query.Order("data ASC").Find(&weatherData).Error
	return weatherData, err
}

func (r *Repository) GetWeatherByDate(cidade, estado string, data time.Time) (*Weather, error) {
	var weather Weather
	err := r.db.Table("previsao_tempo").Where("cidade = ? AND estado = ? AND CAST(data AS DATE) = CAST(? AS DATE)", cidade, estado, data).First(&weather).Error
	if err != nil {
		return nil, err
	}
	return &weather, nil
}
