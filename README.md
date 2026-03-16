# ClimIA Backend

API NestJS que serve previsões climáticas para o [ClimIA](https://www.climia.com.br).

## Arquitetura

```
Frontend (Next.js) → API (NestJS) → DynamoDB (ClimIA-Previsoes)
```

A API lê dados pré-calculados do DynamoDB. Cada cidade do Brasil tem 366 registros (um por dia do ano) com médias históricas de temperatura e precipitação.

## Stack

- **NestJS** — Framework TypeScript
- **AWS DynamoDB** — Banco de dados NoSQL
- **AWS SDK v3** — Client DynamoDB

## Configuração

```bash
npm install
cp .env.example .env
```

### Variáveis de ambiente

| Variável | Descrição | Default |
|----------|-----------|---------|
| `PORT` | Porta da API | `3001` |
| `API_TOKEN` | Token Bearer para autenticação | — |
| `AWS_REGION` | Região AWS | `us-east-1` |
| `DYNAMODB_TABLE` | Nome da tabela DynamoDB | `ClimIA-Previsoes` |

> As credenciais AWS são lidas de `~/.aws/credentials` ou variáveis de ambiente.

## Rodando

```bash
npm run start:dev    # desenvolvimento
npm run build        # build de produção
npm run start:prod   # produção
```

## Endpoints

### `GET /health`
Health check.

### `GET /`
Retorna previsão climática.

| Parâmetro | Obrigatório | Descrição |
|-----------|-------------|-----------|
| `cidade` | Sim | Nome da cidade |
| `estado` | Sim | UF (ex: SP) |
| `data` | * | Data única (YYYY-MM-DD) |
| `datainicio` | * | Data início do range |
| `datafim` | * | Data fim do range |

\* Informe `data` ou `datainicio` + `datafim`.

**Headers:** `Authorization: Bearer <API_TOKEN>`

**Exemplo:**
```
GET /?cidade=Guarulhos&estado=SP&data=2026-03-16
Authorization: Bearer <token>
```

**Resposta:**
```json
[{
  "date": "2026-03-16",
  "city": "Guarulhos",
  "uf": "SP",
  "temp_max": 27.59,
  "temp_min": 19.02,
  "temp_mean": 22.72,
  "precipitation": 7.2,
  "confidence": 0.85
}]
```

## Cache

- **Query cache:** 24h por cidade (todos os 366 dias da cidade)
- **Response cache:** 24h por combinação cidade+datas
