# Serviço de Temperatura por CEP

Este serviço recebe um CEP brasileiro e retorna a temperatura atual em Celsius, Fahrenheit e Kelvin.

## Requisitos

- Go 1.21 ou superior
- Docker e Docker Compose
- Chave de API do OpenWeatherMap (obtida em [OpenWeatherMap](https://openweathermap.org/api))

## Configuração

1. Clone o repositório
2. Copie o arquivo `.env.example` para `.env`:
   ```bash
   cp .env.example .env
   ```
3. Adicione sua chave de API do OpenWeatherMap no arquivo `.env`:
   ```
   OPENWEATHER_API_KEY=sua_chave_aqui
   ```

## Executando com Docker Compose

### Iniciar o Serviço

```bash
docker-compose up
```

Para executar em segundo plano:
```bash
docker-compose up -d
```

### Executar Testes

```bash
docker-compose run test
```

### Parar o Serviço

```bash
docker-compose down
```

## Executando Localmente

1. Instale as dependências:
   ```bash
   go mod download
   ```

2. Execute o servidor:
   ```bash
   go run cmd/main.go
   ```

## Testando a API

Você pode testar a API usando o arquivo `api.http` ou fazendo requisições HTTP diretamente:

```bash
curl "http://localhost:8080/temperature?cep=89010-904"
```

Ou, use o serviço implantado no Google Cloud Run:

```bash
curl "https://go-cloud-run-lab-678036928420.us-central1.run.app/temperature?cep=89010-904"
```

### Exemplo de Resposta

```json
{
    "temp_c": 25.5,
    "temp_f": 77.9,
    "temp_k": 298.65
}
```

## Códigos de Resposta

- 200: Sucesso
- 400: CEP inválido
- 404: CEP não encontrado
- 500: Erro interno do servidor

## Estrutura do Projeto

```
.
├── cmd/
│   └── main.go
├── configs/
│   └── config.go
├── internal/
│   ├── handler.go
│   └── service.go
├── .env
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## APIs Utilizadas

- [ViaCEP](https://viacep.com.br/) - Para obter dados de localização a partir do CEP
- [OpenWeatherMap](https://openweathermap.org/api) - Para obter dados de temperatura

## Fórmulas de Conversão

- Celsius para Fahrenheit: °F = (°C × 9/5) + 32
- Celsius para Kelvin: K = °C + 273.15