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
- MySQL 8.0+

### 1. Configurar banco de dados

Crie um banco MySQL chamado `climia` e execute o script SQL:

```sql
CREATE DATABASE climia;
USE climia;

CREATE TABLE `previsao_tempo` (
  `id` int NOT NULL AUTO_INCREMENT,
  `data` date DEFAULT NULL,
  `temperatura_minima` float DEFAULT NULL,
  `temperatura_maxima` float DEFAULT NULL,
  `temperatura_media` float DEFAULT NULL,
  `cidade` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `estado` varchar(2) COLLATE utf8mb4_unicode_ci NOT NULL,
  `precipitacao` decimal(8,2) NOT NULL DEFAULT '0.00',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_cidade_estado` (`cidade`,`estado`),
  KEY `idx_cidade_estado_data` (`cidade`,`estado`,`data`),
  KEY `idx_data` (`data`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 2. Configurar variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=senha123
DB_NAME=climia
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

- [ ] Migração para Aurora
- [ ] Cache Redis
- [ ] Métricas de performance
- [ ] Documentação OpenAPI
