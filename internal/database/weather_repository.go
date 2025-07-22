package database

import (
	"database/sql"
)

type WeatherRepository interface {
	GetDB() *sql.DB
}

type WeatherRepositoryImpl struct {
	db *sql.DB
}

func NewWeatherRepository(db *sql.DB) WeatherRepository {
	return &WeatherRepositoryImpl{db: db}
}

func (r *WeatherRepositoryImpl) GetDB() *sql.DB {
	return r.db
} 