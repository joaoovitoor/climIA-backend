package services

import (
	"database/sql"
	"fmt"
	"math"
	"sort"
	"time"

	"climia-backend/internal/database"
	"climia-backend/internal/models"
)

type WeatherService struct {
	weatherRepo database.WeatherRepository
}

func NewWeatherService(repo database.WeatherRepository) *WeatherService {
	return &WeatherService{
		weatherRepo: repo,
	}
}

func (s *WeatherService) CalculateForecast(req models.WeatherRequest) ([]models.WeatherResponse, error) {

	if req.Data != "" {
		data, err := time.Parse("2006-01-02", req.Data)
		if err != nil {
			return nil, fmt.Errorf("formato de data inválido: %v", err)
		}

		forecast, err := s.calculateForecastForDate(req.Cidade, req.Estado, data)
		if err != nil {
			return nil, err
		}

		return []models.WeatherResponse{forecast}, nil
	}

	if req.DataInicio != "" && req.DataFim != "" {
		dataInicio, err := time.Parse("2006-01-02", req.DataInicio)
		if err != nil {
			return nil, fmt.Errorf("formato de data início inválido: %v", err)
		}

		dataFim, err := time.Parse("2006-01-02", req.DataFim)
		if err != nil {
			return nil, fmt.Errorf("formato de data fim inválido: %v", err)
		}

		return s.calculateForecastForPeriod(req.Cidade, req.Estado, dataInicio, dataFim)
	}

	return nil, fmt.Errorf("deve fornecer data ou intervalo de datas")
}

func (s *WeatherService) calculateForecastForDate(cidade, estado string, data time.Time) (models.WeatherResponse, error) {
	dia := data.Day()
	mes := int(data.Month())
	ano := data.Year()

	dadosHistoricos, err := s.buscarDadosHistoricosParaTendencia(cidade, estado)
	if err != nil {
		return models.WeatherResponse{}, fmt.Errorf("erro ao buscar dados históricos: %v", err)
	}

	dadosDia := filtrarDadosPorDiaMes(dadosHistoricos, dia, mes)
	if len(dadosDia) == 0 {
		return models.WeatherResponse{}, fmt.Errorf("não há dados históricos para %s/%s no dia %d/%d", cidade, estado, dia, mes)
	}

	previsao := s.calcularPrevisaoInteligente(dadosDia, ano)
	if previsao == nil {
		return models.WeatherResponse{}, fmt.Errorf("erro ao calcular previsão para %s/%s", cidade, estado)
	}

	return models.WeatherResponse{
		Data:              data.Format("2006-01-02"),
		TemperaturaMinima: previsao["minima"],
		TemperaturaMaxima: previsao["maxima"],
		TemperaturaMedia:  previsao["media"],
		Cidade:            cidade,
		Estado:            estado,
		Precipitacao:      previsao["precipitacao"],
	}, nil
}

func (s *WeatherService) calculateForecastForPeriod(cidade, estado string, dataInicio, dataFim time.Time) ([]models.WeatherResponse, error) {
	var previsoes []models.WeatherResponse

	// Buscar dados históricos para calcular tendência
	dadosHistoricos, err := s.buscarDadosHistoricosParaTendencia(cidade, estado)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar dados históricos: %v", err)
	}

	if len(dadosHistoricos) == 0 {
		return nil, fmt.Errorf("não há dados históricos para %s/%s", cidade, estado)
	}

	// Gerar previsões para cada dia do período
	for data := dataInicio; !data.After(dataFim); data = data.AddDate(0, 0, 1) {
		dia := data.Day()
		mes := int(data.Month())
		ano := data.Year()

		// Filtrar dados históricos para este dia/mês
		dadosDia := filtrarDadosPorDiaMes(dadosHistoricos, dia, mes)
		if len(dadosDia) == 0 {
			continue
		}

		// Calcular previsão inteligente com tendência
		previsao := s.calcularPrevisaoInteligente(dadosDia, ano)
		if previsao == nil {
			continue
		}

		previsaoResponse := models.WeatherResponse{
			Data:              data.Format("2006-01-02"),
			TemperaturaMinima: previsao["minima"],
			TemperaturaMaxima: previsao["maxima"],
			TemperaturaMedia:  previsao["media"],
			Cidade:            cidade,
			Estado:            estado,
			Precipitacao:      previsao["precipitacao"],
		}

		previsoes = append(previsoes, previsaoResponse)
	}

	if len(previsoes) == 0 {
		return nil, fmt.Errorf("não há dados históricos para %s/%s no intervalo especificado", cidade, estado)
	}

	return previsoes, nil
}

func (s *WeatherService) buscarDadosHistoricosParaTendencia(cidade, estado string) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			EXTRACT(DAY FROM data) as dia,
			EXTRACT(MONTH FROM data) as mes,
			EXTRACT(YEAR FROM data) as ano,
			AVG(temperatura_minima) as media_minima,
			AVG(temperatura_media) as media_media,
			AVG(temperatura_maxima) as media_maxima,
			AVG(precipitacao) as media_precipitacao
		FROM previsao_tempo 
		WHERE cidade = ? 
			AND estado = ?
			AND data >= DATE_SUB(CURDATE(), INTERVAL 5 YEAR)
			AND data < CURDATE()
		GROUP BY EXTRACT(DAY FROM data), EXTRACT(MONTH FROM data), EXTRACT(YEAR FROM data)
		ORDER BY EXTRACT(YEAR FROM data), EXTRACT(MONTH FROM data), EXTRACT(DAY FROM data)
	`

	rows, err := s.weatherRepo.GetDB().Query(query, cidade, estado)
	if err != nil {
		return nil, fmt.Errorf("erro na query: %v", err)
	}
	defer rows.Close()

	var dados []map[string]interface{}
	for rows.Next() {
		var dia, mes, ano int
		var mediaMinima, mediaMedia, mediaMaxima, mediaPrecipitacao sql.NullFloat64

		err := rows.Scan(&dia, &mes, &ano, &mediaMinima, &mediaMedia, &mediaMaxima, &mediaPrecipitacao)
		if err != nil {
			return nil, fmt.Errorf("erro ao scan: %v", err)
		}

		dado := map[string]interface{}{
			"dia":                  dia,
			"mes":                  mes,
			"ano":                  ano,
			"media_minima":         mediaMinima.Float64,
			"media_media":          mediaMedia.Float64,
			"media_maxima":         mediaMaxima.Float64,
			"media_precipitacao":   mediaPrecipitacao.Float64,
		}
		dados = append(dados, dado)
	}

	return dados, nil
}

func filtrarDadosPorDiaMes(dados []map[string]interface{}, dia, mes int) []map[string]interface{} {
	var dadosFiltrados []map[string]interface{}
	
	for _, dado := range dados {
		if dado["dia"].(int) == dia && dado["mes"].(int) == mes {
			dadosFiltrados = append(dadosFiltrados, dado)
		}
	}
	
	return dadosFiltrados
}

func (s *WeatherService) calcularPrevisaoInteligente(dadosHistoricos []map[string]interface{}, anoAtual int) map[string]float64 {
	if len(dadosHistoricos) == 0 {
		return nil
	}

	// Ordenar dados por ano
	sort.Slice(dadosHistoricos, func(i, j int) bool {
		return dadosHistoricos[i]["ano"].(int) < dadosHistoricos[j]["ano"].(int)
	})

	// Calcular tendências
	tendenciaMinima := s.calcularTendencia(dadosHistoricos, "media_minima")
	tendenciaMedia := s.calcularTendencia(dadosHistoricos, "media_media")
	tendenciaMaxima := s.calcularTendencia(dadosHistoricos, "media_maxima")

	// Pegar último dado
	ultimoDado := dadosHistoricos[len(dadosHistoricos)-1]
	ultimoAno := ultimoDado["ano"].(int)
	anosDesdeUltimo := anoAtual - ultimoAno

	// Aplicar tendência
	previsaoMinima := s.aplicarTendencia(ultimoDado["media_minima"].(float64), tendenciaMinima, anosDesdeUltimo)
	previsaoMedia := s.aplicarTendencia(ultimoDado["media_media"].(float64), tendenciaMedia, anosDesdeUltimo)
	previsaoMaxima := s.aplicarTendencia(ultimoDado["media_maxima"].(float64), tendenciaMaxima, anosDesdeUltimo)

	// Calcular média de precipitação dos últimos 3 anos
	previsaoPrecipitacao := s.calcularMediaPrecipitacao(dadosHistoricos)

	return map[string]float64{
		"minima":       s.formatarTemperatura(previsaoMinima),
		"media":        s.formatarTemperatura(previsaoMedia),
		"maxima":       s.formatarTemperatura(previsaoMaxima),
		"precipitacao": s.formatarPrecipitacao(previsaoPrecipitacao),
	}
}

func (s *WeatherService) calcularTendencia(dados []map[string]interface{}, campo string) float64 {
	if len(dados) < 2 {
		return 0
	}

	n := len(dados)
	var sumX, sumY, sumXY, sumX2 float64

	for i, dado := range dados {
		x := float64(i)
		y := dado[campo].(float64)

		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	denominador := float64(n)*sumX2 - sumX*sumX
	if denominador == 0 {
		return 0
	}

	tendencia := (float64(n)*sumXY - sumX*sumY) / denominador
	return tendencia
}

func (s *WeatherService) aplicarTendencia(valorBase, tendencia float64, anos int) float64 {
	fatorDecaimento := math.Exp(-float64(anos) * 0.1)
	ajuste := tendencia * float64(anos) * fatorDecaimento

	return valorBase + ajuste
}

func (s *WeatherService) calcularMediaPrecipitacao(dadosOrdenados []map[string]interface{}) float64 {
	ultimos3Anos := dadosOrdenados
	if len(dadosOrdenados) > 3 {
		ultimos3Anos = dadosOrdenados[len(dadosOrdenados)-3:]
	}

	if len(ultimos3Anos) == 0 {
		return 0
	}

	var somaPrecipitacao float64
	for _, dado := range ultimos3Anos {
		somaPrecipitacao += dado["media_precipitacao"].(float64)
	}

	return somaPrecipitacao / float64(len(ultimos3Anos))
}

func (s *WeatherService) formatarTemperatura(valor float64) float64 {
	return math.Round(valor)
}

func (s *WeatherService) formatarPrecipitacao(valor float64) float64 {
	return math.Round(valor)
} 