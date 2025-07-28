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

	sql := "SELECT * FROM previsao_tempo WHERE cidade = '" + cidade + "' AND estado = '" + estado + "'"

	if dataInicio != nil && dataFim != nil {
		sql += " AND data >= '" + dataInicio.Format("2006-01-02") + "' AND data <= '" + dataFim.Format("2006-01-02") + "'"
	} else if dataInicio != nil {
		sql += " AND data >= '" + dataInicio.Format("2006-01-02") + "'"
	} else if dataFim != nil {
		sql += " AND data <= '" + dataFim.Format("2006-01-02") + "'"
	}

	sql += " ORDER BY data ASC"

	err := r.db.Raw(sql).Scan(&weatherData).Error
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
