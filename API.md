# ClimIA API - Documentação

## Endpoints

### GET /health
Verifica se a API está funcionando.

**Resposta:**
```json
{
  "status": "ok",
  "message": "ClimIA API is running"
}
```

### GET /
Calcula previsão do tempo.

**Parâmetros:**
- `cidade` (obrigatório): Nome da cidade
- `estado` (obrigatório): Sigla do estado
- `data` (opcional): Data específica (YYYY-MM-DD)
- `datainicio` (opcional): Data inicial do período (YYYY-MM-DD)
- `datafim` (opcional): Data final do período (YYYY-MM-DD)

**Exemplos:**

**Previsão para uma data específica:**
```
GET /?cidade=Guarulhos&estado=SP&data=2025-11-01
```

**Previsão para um período:**
```
GET /?cidade=Guarulhos&estado=SP&datainicio=2025-11-01&datafim=2025-11-07
```

**Resposta (data única):**
```json
{
  "cidade": "Guarulhos",
  "estado": "SP",
  "data": "2025-11-01",
  "temperatura_minima": 16,
  "temperatura_maxima": 30,
  "temperatura_media": 22,
  "precipitacao": 4
}
```

**Resposta (período):**
```json
[
  {
    "cidade": "Guarulhos",
    "estado": "SP",
    "data": "2025-11-01",
    "temperatura_minima": 16,
    "temperatura_maxima": 30,
    "temperatura_media": 22,
    "precipitacao": 4
  },
  {
    "cidade": "Guarulhos",
    "estado": "SP",
    "data": "2025-11-02",
    "temperatura_minima": 18,
    "temperatura_maxima": 32,
    "temperatura_media": 24,
    "precipitacao": 2
  }
]
```

## URLs

**Local:**
- http://localhost:8080

**AWS Lambda:**
- https://hls852t472.execute-api.us-east-1.amazonaws.com/prod/

## Códigos de Erro

- `400`: Parâmetros inválidos ou obrigatórios não fornecidos
- `500`: Erro interno do servidor 