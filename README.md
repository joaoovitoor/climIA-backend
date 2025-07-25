# ClimIA Backend

API REST em Go para previsão climática baseada em dados históricos, com estrutura simplificada e otimizada.

## 🏗️ Arquitetura

O projeto segue uma estrutura modular organizada:

- **configs/**: Configurações da aplicação
- **internal/modules/**: Módulos da aplicação
  - **weather/**: Módulo de previsão do tempo
    - **handler.go**: Controllers HTTP
    - **service.go**: Lógica de negócio
    - **repository.go**: Camada de dados
    - **model.go**: Entidades e DTOs
- **pkg/database/**: Configuração do banco de dados
- **cmd/**: Pontos de entrada da aplicação

## 🚀 Como executar

### Pré-requisitos

- Go 1.21+
- Aurora PostgreSQL (instância de leitura)

### 1. Configurar banco de dados

Crie um banco PostgreSQL chamado `postgres` e execute o script SQL:

```sql
CREATE TABLE previsao_tempo (
  id SERIAL PRIMARY KEY,
  data DATE,
  temperatura_minima FLOAT,
  temperatura_maxima FLOAT,
  temperatura_media FLOAT,
  cidade VARCHAR(255) NOT NULL,
  estado VARCHAR(2) NOT NULL,
  precipitacao DECIMAL(8,2) NOT NULL DEFAULT 0.00,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_cidade_estado ON previsao_tempo(cidade, estado);
CREATE INDEX idx_cidade_estado_data ON previsao_tempo(cidade, estado, data);
CREATE INDEX idx_data ON previsao_tempo(data);
```

### 2. Configurar variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto:

```env
DB_CONNECTION_STRING=postgres://username:password@aurora-instance:5432/database?sslmode=require
PORT=8080
ENV=development
API_TOKEN=your_api_token_here
```

**Nota**: O arquivo `.env` não é versionado no Git por questões de segurança.

### 3. Instalar dependências

```bash
go mod tidy
```

### 4. Executar a aplicação

```bash
go run cmd/api/main.go
```

A API estará disponível em `http://localhost:8080`

## 📡 Endpoints

### GET /

Gera previsão climática baseada em dados históricos.

#### Parâmetros:

- `cidade` (obrigatório): Nome da cidade
- `estado` (obrigatório): Sigla do estado (ex: SP, RJ)
- `data` (opcional): Data específica no formato YYYY-MM-DD
- `datainicio` (opcional): Data de início do intervalo
- `datafim` (opcional): Data de fim do intervalo

#### Exemplos de uso:

```bash
# Previsão para data específica
curl -H "Authorization: Bearer YOUR_API_TOKEN" "http://localhost:8080/?cidade=Guarulhos&estado=SP&data=2025-11-01"

# Previsão para intervalo de datas
curl -H "Authorization: Bearer YOUR_API_TOKEN" "http://localhost:8080/?cidade=Guarulhos&estado=SP&datainicio=2025-11-01&datafim=2025-11-07"

# Informações da API
curl -H "Authorization: Bearer YOUR_API_TOKEN" "http://localhost:8080/"
```

### GET /health

Verifica o status da API.

```bash
curl "http://localhost:8080/health"
```

#### Resposta:

```json
{
  "status": "ok",
  "message": "ClimIA API is running"
}
```

## 🔧 Desenvolvimento

### Hot Reload

Para desenvolvimento com hot reload, use o Air:

```bash
air
```

### Build

Para gerar o executável:

```bash
go build -o main cmd/api/main.go
```

## 🔐 Autenticação

A API utiliza autenticação via Bearer Token. Inclua o header `Authorization: Bearer YOUR_API_TOKEN` em todas as requisições.

## 🚀 Deploy

O projeto está configurado para deploy automático na AWS Lambda via GitHub Actions:

- **Trigger**: Push para branch `main`
- **Infraestrutura**: AWS Lambda + API Gateway
- **Banco**: Aurora PostgreSQL (instância de leitura)
- **Variáveis**: Configuradas via GitHub Secrets

## 🏛️ Estrutura do Projeto

```
climIA-backend/
├── cmd/
│   └── api/main.go          # Ponto de entrada da API
├── configs/
│   └── config.go            # Configurações da aplicação
├── internal/
│   ├── app/
│   │   └── app.go           # Configuração da aplicação
│   └── modules/
│       └── weather/         # Módulo de previsão do tempo
│           ├── handler.go    # Controllers HTTP
│           ├── service.go    # Lógica de negócio
│           ├── repository.go # Camada de dados
│           └── model.go      # Entidades e DTOs
├── pkg/
│   └── database/
│       └── connection.go     # Configuração do banco
├── .github/workflows/
│   └── deploy.yml           # CI/CD para AWS Lambda
├── go.mod
└── go.sum
```
