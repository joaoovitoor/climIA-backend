# ClimIA Backend

API REST em Go para previsão climática baseada em dados históricos, com estrutura simplificada e otimizada.

## 🏗️ Arquitetura

O projeto segue uma estrutura intermediária organizada:

- **models/**: Entidades e DTOs
- **services/**: Lógica de negócio
- **handlers/**: Controllers HTTP
- **database/**: Camada de dados e repositórios
- **routes/**: Definição de rotas
- **server/**: Configuração do servidor

## 🚀 Como executar

### Pré-requisitos

- Go 1.21+
- PostgreSQL 12+ ou Aurora PostgreSQL

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
DB_CONNECTION_STRING=postgres://username:password@host:5432/database?sslmode=require
PORT=8080
ENV=development
API_TOKEN=your_api_token_here
```

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
curl "http://localhost:8080/?cidade=Guarulhos&estado=SP&data=2025-11-01"

# Previsão para intervalo de datas
curl "http://localhost:8080/?cidade=Guarulhos&estado=SP&datainicio=2025-11-01&datafim=2025-11-07"

# Informações da API
curl "http://localhost:8080/"
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

## 📊 Performance

- **Latência**: ~44ms para 7 dias de previsão
- **Tamanho**: ~119KB por resposta
- **Algoritmo**: Cálculo inteligente com tendência e decaimento

## 🏛️ Estrutura do Projeto

```
climIA-backend/
├── cmd/api/main.go           # Ponto de entrada
├── config/config.go          # Configurações do banco
├── internal/
│   ├── models/              # Entidades e DTOs
│   ├── services/            # Lógica de negócio
│   ├── handlers/            # Controllers HTTP
│   ├── routes/              # Definição de rotas
│   ├── database/            # Camada de dados
│   ├── server/              # Servidor HTTP
│   └── config/              # Configurações da app
├── go.mod
└── go.sum
```

## 🚀 Próximos Passos

- [x] Migração para Aurora PostgreSQL
- [ ] Cache Redis
- [ ] Métricas de performance
- [ ] Documentação OpenAPI
