package weather

import (
	"fmt"
	"time"

	"climia-backend/configs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBRepository struct {
	svc    *dynamodb.DynamoDB
	config *configs.Config
}

type DynamoDBWeather struct {
	PK                string  `json:"PK"`
	SK                string  `json:"SK"`
	Cidade            string  `json:"cidade"`
	Estado            string  `json:"estado"`
	Dia               int     `json:"dia"`
	Mes               int     `json:"mes"`
	TemperaturaMinima float64 `json:"temperatura_minima"`
	TemperaturaMaxima float64 `json:"temperatura_maxima"`
	TemperaturaMedia  float64 `json:"temperatura_media"`
	Precipitacao      float64 `json:"precipitacao"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

func NewDynamoDBRepository(config *configs.Config) (*DynamoDBRepository, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(config.DynamoAccessKey, config.DynamoSecret, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao criar sessão AWS: %v", err)
	}

	svc := dynamodb.New(sess)
	return &DynamoDBRepository{svc: svc, config: config}, nil
}

func (r *DynamoDBRepository) GetWeatherByDate(cidade, estado string, data time.Time) (*WeatherResponse, error) {
	dia := data.Day()
	mes := int(data.Month())

	pk := fmt.Sprintf("CITY#%s#%s", cidade, estado)
	sk := fmt.Sprintf("DATE#%02d#%02d", dia, mes)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.config.DynamoTableName),
		KeyConditionExpression: aws.String("PK = :pk AND SK = :sk"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(pk),
			},
			":sk": {
				S: aws.String(sk),
			},
		},
	}

	result, err := r.svc.Query(input)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar DynamoDB: %v", err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("dados não encontrados para %s/%s na data %s", cidade, estado, data.Format("2006-01-02"))
	}

	var dynamoWeather DynamoDBWeather
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &dynamoWeather)
	if err != nil {
		return nil, fmt.Errorf("erro ao deserializar dados: %v", err)
	}

	return &WeatherResponse{
		Cidade:            dynamoWeather.Cidade,
		Estado:            dynamoWeather.Estado,
		Data:              data.Format("2006-01-02"),
		TemperaturaMinima: dynamoWeather.TemperaturaMinima,
		TemperaturaMaxima: dynamoWeather.TemperaturaMaxima,
		TemperaturaMedia:  dynamoWeather.TemperaturaMedia,
		Precipitacao:      dynamoWeather.Precipitacao,
	}, nil
}

func (r *DynamoDBRepository) GetWeatherDataForPeriod(cidade, estado string, dataInicio, dataFim time.Time) ([]WeatherResponse, error) {
	pk := fmt.Sprintf("CITY#%s#%s", cidade, estado)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.config.DynamoTableName),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(pk),
			},
		},
	}

	result, err := r.svc.Query(input)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar DynamoDB: %v", err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("dados não encontrados para %s/%s", cidade, estado)
	}

	var responses []WeatherResponse
	for _, item := range result.Items {
		var dynamoWeather DynamoDBWeather
		err := dynamodbattribute.UnmarshalMap(item, &dynamoWeather)
		if err != nil {
			continue
		}

		// Reconstruir a data a partir do dia e mês (usando o ano da data de início)
		ano := dataInicio.Year()
		data := time.Date(ano, time.Month(dynamoWeather.Mes), dynamoWeather.Dia, 0, 0, 0, 0, time.UTC)

		// Verificar se a data está no período solicitado
		if (data.After(dataInicio) || data.Equal(dataInicio)) && (data.Before(dataFim) || data.Equal(dataFim)) {
			responses = append(responses, WeatherResponse{
				Cidade:            dynamoWeather.Cidade,
				Estado:            dynamoWeather.Estado,
				Data:              data.Format("2006-01-02"),
				TemperaturaMinima: dynamoWeather.TemperaturaMinima,
				TemperaturaMaxima: dynamoWeather.TemperaturaMaxima,
				TemperaturaMedia:  dynamoWeather.TemperaturaMedia,
				Precipitacao:      dynamoWeather.Precipitacao,
			})
		}
	}

	if len(responses) == 0 {
		return nil, fmt.Errorf("nenhum dado encontrado para %s/%s no período especificado", cidade, estado)
	}

	return responses, nil
}

func (r *DynamoDBRepository) GetAllCitiesByState(estado string) ([]string, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String(r.config.DynamoTableName),
		FilterExpression: aws.String("estado = :estado"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":estado": {
				S: aws.String(estado),
			},
		},
	}

	result, err := r.svc.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar DynamoDB: %v", err)
	}

	cities := make(map[string]bool)
	for _, item := range result.Items {
		if cidade, ok := item["cidade"]; ok && cidade.S != nil {
			cities[*cidade.S] = true
		}
	}

	var cityList []string
	for cidade := range cities {
		cityList = append(cityList, cidade)
	}

	return cityList, nil
}
