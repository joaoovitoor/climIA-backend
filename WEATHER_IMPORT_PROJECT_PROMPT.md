# PROMPT: Sistema de Importação de Dados Meteorológicos - LOCAL

## Objetivo

Criar um sistema **LOCAL** em Go para importar dados meteorológicos da API Open-Meteo para um banco de dados MySQL. O sistema roda apenas localmente, sem Lambda ou cloud.

## Estrutura do Projeto

### Tecnologias

- **Linguagem**: Go 1.21+
- **Banco de Dados**: MySQL 8.0+ (local)
- **API Externa**: Open-Meteo (https://open-meteo.com/)
- **Execução**: Apenas local (não Lambda, não cloud)
- **Logging**: Logs locais simples
- **Configuração**: Arquivo .env

### Estrutura de Diretórios

```
weather-import/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── connection.go
│   │   ├── city_repository.go
│   │   └── weather_repository.go
│   ├── models/
│   │   ├── city.go
│   │   ├── weather.go
│   │   └── weather_api.go
│   └── services/
│       └── weather_import_service.go
├── scripts/
│   ├── database_schema.sql
│   └── seed_cities.sql
├── config/
│   └── cities.json
├── .env
├── go.mod
├── go.sum
└── README.md
```

## Configuração

### Variáveis de Ambiente (.env)

```env
# Database LOCAL
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=senha123
DB_NAME=climia
DB_CONNECTION_STRING=root:senha123@tcp(localhost:3306)/climia?parseTime=true

# Open-Meteo API
OPEN_METEO_API_KEY=Iy7gvUmYUUQqMs2j
OPEN_METEO_BASE_URL=https://customer-archive-api.open-meteo.com/v1/archive
OPEN_METEO_TIMEOUT=30s

# Application
APP_ENV=development
LOG_LEVEL=info

# Import Settings
IMPORT_DAYS_BACK=7
IMPORT_DELAY_BETWEEN_REQUESTS=100ms
```

## Banco de Dados

### Tabela: cities

```sql
CREATE TABLE cities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    estado VARCHAR(50) NOT NULL,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    ultima_atualizacao DATETIME NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_cidade_estado (nome, estado),
    INDEX idx_ultima_atualizacao (ultima_atualizacao)
);
```

### Tabela: weather

```sql
CREATE TABLE weather (
    id INT AUTO_INCREMENT PRIMARY KEY,
    data DATE NOT NULL,
    temperatura_minima DECIMAL(5, 2) NOT NULL,
    temperatura_maxima DECIMAL(5, 2) NOT NULL,
    temperatura_media DECIMAL(5, 2) NOT NULL,
    cidade VARCHAR(100) NOT NULL,
    estado VARCHAR(50) NOT NULL,
    precipitacao DECIMAL(8, 2) NOT NULL DEFAULT 0.00,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_cidade_estado_data (cidade, estado, data),
    INDEX idx_cidade_estado (cidade, estado),
    INDEX idx_data (data)
);
```

## Modelos de Dados

### City Model

```go
type City struct {
    ID                uint      `json:"id"`
    Nome              string    `json:"nome"`
    Estado            string    `json:"estado"`
    Latitude          float64   `json:"latitude"`
    Longitude         float64   `json:"longitude"`
    UltimaAtualizacao time.Time `json:"ultima_atualizacao"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
}
```

### Weather Model

```go
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
```

### API Response Models

```go
type OpenMeteoAPIResponse struct {
    Daily OpenMeteoDaily `json:"daily"`
}

type OpenMeteoDaily struct {
    Time                []string  `json:"time"`
    Temperature2mMax    []float64 `json:"temperature_2m_max"`
    Temperature2mMin    []float64 `json:"temperature_2m_min"`
    Temperature2mMean   []float64 `json:"temperature_2m_mean"`
    PrecipitationSum    []float64 `json:"precipitation_sum"`
}

type WeatherData struct {
    Data              string  `json:"data"`
    TemperaturaMinima float64 `json:"temperatura_minima"`
    TemperaturaMaxima float64 `json:"temperatura_maxima"`
    TemperaturaMedia  float64 `json:"temperatura_media"`
    Precipitacao      float64 `json:"precipitacao"`
}
```

## Funcionalidades Principais

### 1. Importação de Dados

- **Importação por cidade**: Importar dados para uma cidade específica
- **Importação em lote**: Importar dados para todas as cidades
- **Controle de duplicatas**: Remove dados existentes antes de inserir novos
- **Retry mechanism**: Tentativas automáticas em caso de falha
- **Rate limiting**: Controle de velocidade das requisições

### 2. Execução Local

- **Script simples**: Execução manual via linha de comando
- **Logs locais**: Logs simples no console
- **Configuração local**: Arquivo .env para configurações

### 3. Monitoramento Simples

- **Logs no console**: Logs detalhados durante execução
- **Status de execução**: Mostra sucesso/falha por cidade
- **Resumo final**: Relatório final da importação

## Arquivo de Cidades (config/cities.json)

```json
{
  "cities": [
    {
      "nome": "São Paulo",
      "estado": "SP",
      "latitude": -23.5505,
      "longitude": -46.6333
    },
    {
      "nome": "Rio de Janeiro",
      "estado": "RJ",
      "latitude": -22.9068,
      "longitude": -43.1729
    },
    {
      "nome": "Belo Horizonte",
      "estado": "MG",
      "latitude": -19.9167,
      "longitude": -43.9345
    },
    {
      "nome": "Salvador",
      "estado": "BA",
      "latitude": -12.9714,
      "longitude": -38.5011
    },
    {
      "nome": "Fortaleza",
      "estado": "CE",
      "latitude": -3.7319,
      "longitude": -38.5267
    },
    {
      "nome": "Brasília",
      "estado": "DF",
      "latitude": -15.7942,
      "longitude": -47.8822
    },
    {
      "nome": "Curitiba",
      "estado": "PR",
      "latitude": -25.4284,
      "longitude": -49.2733
    },
    {
      "nome": "Recife",
      "estado": "PE",
      "latitude": -8.0476,
      "longitude": -34.877
    },
    {
      "nome": "Porto Alegre",
      "estado": "RS",
      "latitude": -30.0346,
      "longitude": -51.2177
    },
    {
      "nome": "Manaus",
      "estado": "AM",
      "latitude": -3.119,
      "longitude": -60.0217
    }
  ]
}
```

## Dependências (go.mod)

```go
module weather-import

go 1.21

require (
    github.com/go-sql-driver/mysql v1.7.0
    github.com/joho/godotenv v1.5.1
    github.com/sirupsen/logrus v1.9.3
)
```

## Como Executar

### Setup Inicial

```bash
# 1. Criar banco de dados
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS climia;"

# 2. Executar schema
mysql -u root -p climia < scripts/database_schema.sql

# 3. Inserir cidades
mysql -u root -p climia < scripts/seed_cities.sql

# 4. Configurar .env
cp .env.example .env
# Editar .env com suas configurações

# 5. Build e executar
go build -o weather-import cmd/main.go
./weather-import
```

### Execução Manual

```bash
# Importar dados para todas as cidades
./weather-import --import-all

# Importar dados para cidade específica
./weather-import --city "São Paulo" --state "SP"

# Verificar status das cidades
./weather-import --status
```

## Funcionalidades do Sistema

### 1. Importação Automática

- Busca cidades que não foram atualizadas nos últimos 7 dias
- Importa dados dos últimos 7 dias por padrão
- Remove dados duplicados antes de inserir novos

### 2. Logs Detalhados

- Mostra progresso por cidade
- Exibe erros específicos por cidade
- Relatório final com estatísticas

### 3. Configuração Flexível

- Arquivo .env para configurações
- Arquivo JSON para lista de cidades
- Parâmetros via linha de comando

## Estrutura do Código

### main.go

```go
func main() {
    // Carregar configuração
    config := config.LoadConfig()

    // Conectar ao banco
    db := database.NewConnection(config)

    // Criar serviços
    weatherService := services.NewWeatherImportService(config, db)

    // Executar importação baseada nos argumentos
    weatherService.Run()
}
```

### weather_import_service.go

```go
type WeatherImportService struct {
    config *config.Config
    db     *sql.DB
    logger *logrus.Logger
}

func (s *WeatherImportService) ImportAllCities() error
func (s *WeatherImportService) ImportCity(cidade, estado string) error
func (s *WeatherImportService) CheckStatus() error
```

## Cronograma de Desenvolvimento

### Fase 1: Estrutura Base

- [ ] Setup do projeto Go
- [ ] Configuração de banco de dados local
- [ ] Modelos de dados
- [ ] Repositórios básicos
- [ ] Configuração de logging

### Fase 2: Integração com API

- [ ] Cliente HTTP para Open-Meteo
- [ ] Parsing de respostas da API
- [ ] Tratamento de erros
- [ ] Rate limiting
- [ ] Retry mechanism

### Fase 3: Importação de Dados

- [ ] Serviço de importação
- [ ] Controle de duplicatas
- [ ] Processamento em lotes
- [ ] Validação de dados
- [ ] Logs de importação

### Fase 4: Interface de Linha de Comando

- [ ] Argumentos de linha de comando
- [ ] Comandos para diferentes operações
- [ ] Help e documentação
- [ ] Relatórios de execução

## Considerações Finais

- **Execução LOCAL**: Sistema roda apenas na máquina local
- **Simplicidade**: Sem complexidades de cloud/Lambda
- **Flexibilidade**: Configuração via arquivos e linha de comando
- **Confiabilidade**: Tratamento de erros e retry mechanism
- **Monitoramento**: Logs detalhados no console
