# Temperatura por CEP

Sistema em Go que retorna a temperatura atual (Celsius, Fahrenheit e Kelvin) a partir de um CEP brasileiro.

## Endpoints

- `GET /{cep}` - Retorna a temperatura para o CEP informado
- `GET /health` - Health check da aplicação

## Exemplo

```bash
curl http://localhost:8080/01001000
```

```cloudrun
 curl https://temperatura-por-cep-502159538479.europe-west1.run.app/01001000
```
Resposta:

```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.7
}
```

## Execução local

### Com Docker

```bash
# Construir imagem
docker build -t temperatura-por-cep .

# Executar
docker run -p 8080:8080 temperatura-por-cep
```

### Sem Docker

```bash
# Baixar dependências
go mod download

# Executar
go run cmd/main.go
```

## Variáveis de ambiente

| Variável | Descrição                                 | Padrão |
|----------|-------------------------------------------|--------|
| `PORT`   | Porta onde o servidor irá escutar         | `8080` |

## Testes

### Executar todos os testes

```bash
# Com cobertura
go test -cover ./...

# Verbose
go test -v ./...

# Apenas pacote específico
go test ./pkg/converter
go test ./internal/service
go test ./internal/handler
```

### Testes unitários

- **pkg/converter**: Testa funções de conversão de temperatura
- **internal/service**: Testa serviços de integração (ViaCEP e Open-Meteo) com HTTP mock
- **internal/handler**: Testa handlers HTTP com mocks

## Deploy no Google Cloud Run

### Pré-requisitos

- Google Cloud SDK instalado
- Projeto no Google Cloud com Cloud Run API habilitada

### Passos

```bash
# 1. Construir imagem e marcar para GCR
docker build -t gcr.io/PROJECT_ID/temperatura-por-cep .

# 2. Enviar para Container Registry
docker push gcr.io/PROJECT_ID/temperatura-por-cep

# 3. Deploy no Cloud Run
gcloud run deploy temperatura-por-cep \
  --image gcr.io/PROJECT_ID/temperatura-por-cep \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

## APIs externas

- **ViaCEP**: https://viacep.com.br/ - Busca de endereço por CEP (gratuita)
- **Open-Meteo**: https://open-meteo.com/ - Dados meteorológicos (gratuita, sem necessidade de chave API)

## Estrutura do projeto

```
temperaturaPorCep/
├── cmd/
│   └── main.go              # Ponto de entrada da aplicação
├── internal/
│   ├── handler/
│   │   ├── handlers.go      # Handlers HTTP
│   │   └── handlers_test.go # Testes dos handlers
│   ├── service/
│   │   ├── viacep.go        # Serviço ViaCEP
│   │   ├── viacep_test.go
│   │   ├── weather.go       # Serviço Open-Meteo (clima)
│   │   └── weather_test.go
│   └── model/
│       └── models.go        # Estruturas de dados
├── pkg/
│   └── converter/
│       ├── converter.go     # Funções de conversão
│       └── converter_test.go
├── Dockerfile               # Multi-stage build
├── .dockerignore
├── go.mod
└── README.md
```

## Formato de CEP

O sistema aceita CEPs no formato de 8 dígitos (sem hífen). Exemplo: `01001000`

## Códigos HTTP

| Código | Significado           | Mensagem                  |
|--------|----------------------|---------------------------|
| 200    | OK                   | `{temp_C, temp_F, temp_K}` |
| 404    | CEP não encontrado   | `can not find zipcode`    |
| 422    | Formato inválido     | `invalid zipcode`         |
| 500    | Erro interno         | `internal server error`   |

## Conversões

- Celsius para Fahrenheit: `F = C * 1.8 + 32`
- Celsius para Kelvin: `K = C + 273.15`

As temperaturas são arredondadas para 1 casa decimal na resposta.

## Tempo real

A aplicação consulta dados meteorológicos em tempo real através da API Open-Meteo, que não requer chave de autenticação. Os dados são atualizados continuamente a partir dos melhores modelos meteorológicos globais.

## Licença

MIT
