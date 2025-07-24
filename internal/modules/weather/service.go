package weather

import (
	"database/sql"
	"fmt"
	"math"
	"sort"
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

	if req.Data != "" {
		data, err := time.Parse("2006-01-02", req.Data)
		if err != nil {
			return nil, fmt.Errorf("formato de data inválido. Use YYYY-MM-DD")
		}

		forecast, err := s.getWeatherDataForDate(req.Cidade, req.Estado, data)
		if err != nil {
			return nil, err
		}

		return []WeatherResponse{forecast}, nil
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

		return s.getWeatherDataForPeriod(req.Cidade, req.Estado, dataInicio, dataFim)
	}

	return nil, fmt.Errorf("deve fornecer data ou intervalo de datas")
}

func (s *Service) getWeatherDataForDate(cidade, estado string, data time.Time) (WeatherResponse, error) {
	hoje := time.Now().Truncate(24 * time.Hour)
	dataTruncada := data.Truncate(24 * time.Hour)
	
	if dataTruncada.Before(hoje) {
		weather, err := s.repository.GetWeatherByDate(cidade, estado, data)
		if err != nil {
			return WeatherResponse{}, fmt.Errorf("dados não encontrados para a data especificada")
		}
		
		return WeatherResponse{
			Data:              weather.Data.Format("2006-01-02"),
			TemperaturaMinima: weather.TemperaturaMinima,
			TemperaturaMaxima: weather.TemperaturaMaxima,
			TemperaturaMedia:  weather.TemperaturaMedia,
			Cidade:            weather.Cidade,
			Estado:            weather.Estado,
			Precipitacao:      weather.Precipitacao,
		}, nil
	} else {
		return s.calculateForecastForDate(cidade, estado, data)
	}
}

func (s *Service) calculateForecastForDate(cidade, estado string, data time.Time) (WeatherResponse, error) {
	dia := data.Day()
	mes := int(data.Month())
	ano := data.Year()

	dadosHistoricos, err := s.buscarDadosHistoricosParaTendencia(cidade, estado)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("erro ao buscar dados históricos: %v", err)
	}

	dadosDia := filtrarDadosPorDiaMes(dadosHistoricos, dia, mes)
	if len(dadosDia) == 0 {
		return WeatherResponse{}, fmt.Errorf("não há dados históricos para %s/%s no dia %d/%d", cidade, estado, dia, mes)
	}

	previsao := s.calcularPrevisaoInteligente(dadosDia, ano)
	if previsao == nil {
		return WeatherResponse{}, fmt.Errorf("erro ao calcular previsão para %s/%s", cidade, estado)
	}

	return WeatherResponse{
		Data:              data.Format("2006-01-02"),
		TemperaturaMinima: previsao["minima"],
		TemperaturaMaxima: previsao["maxima"],
		TemperaturaMedia:  previsao["media"],
		Cidade:            cidade,
		Estado:            estado,
		Precipitacao:      previsao["precipitacao"],
	}, nil
}

func (s *Service) getWeatherDataForPeriod(cidade, estado string, dataInicio, dataFim time.Time) ([]WeatherResponse, error) {
	var resultados []WeatherResponse

	for data := dataInicio; !data.After(dataFim); data = data.AddDate(0, 0, 1) {
		weather, err := s.getWeatherDataForDate(cidade, estado, data)
		if err != nil {
			continue
		}
		resultados = append(resultados, weather)
	}

	if len(resultados) == 0 {
		return nil, fmt.Errorf("nenhum dado encontrado para %s, %s", cidade, estado)
	}

	return resultados, nil
}

func (s *Service) buscarDadosHistoricosParaTendencia(cidade, estado string) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			DAY(data) as dia,
			MONTH(data) as mes,
			YEAR(data) as ano,
			AVG(temperatura_minima) as media_minima,
			AVG(temperatura_media) as media_media,
			AVG(temperatura_maxima) as media_maxima,
			AVG(precipitacao) as media_precipitacao
		FROM previsao_tempo 
		WHERE cidade = ? 
			AND estado = ?
			AND data >= DATE_SUB(CURDATE(), INTERVAL 5 YEAR)
			AND data < CURDATE()
		GROUP BY DAY(data), MONTH(data), YEAR(data)
		ORDER BY YEAR(data), MONTH(data), DAY(data)
	`

	rows, err := s.repository.db.Raw(query, cidade, estado).Rows()
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
			"dia":                dia,
			"mes":                mes,
			"ano":                ano,
			"media_minima":       mediaMinima.Float64,
			"media_media":        mediaMedia.Float64,
			"media_maxima":       mediaMaxima.Float64,
			"media_precipitacao": mediaPrecipitacao.Float64,
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

func (s *Service) calcularPrevisaoInteligente(dadosHistoricos []map[string]interface{}, anoAtual int) map[string]float64 {
	if len(dadosHistoricos) == 0 {
		return nil
	}

	sort.Slice(dadosHistoricos, func(i, j int) bool {
		return dadosHistoricos[i]["ano"].(int) < dadosHistoricos[j]["ano"].(int)
	})

	tendenciaMinima := s.calcularTendencia(dadosHistoricos, "media_minima")
	tendenciaMedia := s.calcularTendencia(dadosHistoricos, "media_media")
	tendenciaMaxima := s.calcularTendencia(dadosHistoricos, "media_maxima")

	ultimoDado := dadosHistoricos[len(dadosHistoricos)-1]
	ultimoAno := ultimoDado["ano"].(int)
	anosDesdeUltimo := anoAtual - ultimoAno

	previsaoMinima := s.aplicarTendencia(ultimoDado["media_minima"].(float64), tendenciaMinima, anosDesdeUltimo)
	previsaoMedia := s.aplicarTendencia(ultimoDado["media_media"].(float64), tendenciaMedia, anosDesdeUltimo)
	previsaoMaxima := s.aplicarTendencia(ultimoDado["media_maxima"].(float64), tendenciaMaxima, anosDesdeUltimo)

	previsaoPrecipitacao := s.calcularMediaPrecipitacao(dadosHistoricos)

	return map[string]float64{
		"minima":       s.formatarTemperatura(previsaoMinima),
		"media":        s.formatarTemperatura(previsaoMedia),
		"maxima":       s.formatarTemperatura(previsaoMaxima),
		"precipitacao": s.formatarPrecipitacao(previsaoPrecipitacao),
	}
}

func (s *Service) calcularTendencia(dados []map[string]interface{}, campo string) float64 {
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

func (s *Service) aplicarTendencia(valorBase, tendencia float64, anos int) float64 {
	fatorDecaimento := math.Exp(-float64(anos) * 0.1)
	ajuste := tendencia * float64(anos) * fatorDecaimento

	return valorBase + ajuste
}

func (s *Service) calcularMediaPrecipitacao(dadosOrdenados []map[string]interface{}) float64 {
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

func (s *Service) formatarTemperatura(valor float64) float64 {
	return math.Round(valor)
}

func (s *Service) formatarPrecipitacao(valor float64) float64 {
	return math.Round(valor)
} 