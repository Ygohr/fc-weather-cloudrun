# Location Weather - ViaCEP + WeatherAPI + CloudRun

## Project Overview

Weather Cloud Run is a standalone Go HTTP API that receives a Brazilian ZIP Code (CEP), resolves its city using ViaCEP, fetches current weather data from WeatherAPI, and returns the temperature in Celsius, Fahrenheit, and Kelvin.

## Cloud Run

Application URL: https://location-weather-164168176181.us-central1.run.app/weather?cep=your-cep

## Architecture

The project follows a simple Clean Architecture organization:

- `cmd/api`: application entry point.
- `infrastructure/http`: HTTP handler, router, and server setup.
- `internal/domain`: domain models, validation, and domain errors.
- `internal/usecase`: application use case orchestration.
- `internal/service`: external API clients for ViaCEP and WeatherAPI.
- `internal/util`: temperature conversion helpers.
- `tests`: reusable test helpers and HTTP mocks.

## Technologies

- Go 1.24+
- Standard `net/http`
- `context`
- Go modules
- Docker
- Docker Compose
- Testify
- `httptest`

## Requirements

- Go 1.24 or newer
- Docker
- Docker Compose
- WeatherAPI key

## Environment Variables

- `PORT`: HTTP server port. Default: `8080`.
- `WEATHER_API_KEY`: required WeatherAPI key.

Create a local `.env` file if you want Docker Compose to load values automatically:

```env
PORT=8080
WEATHER_API_KEY=your-api-key
```

If `WEATHER_API_KEY` is missing, the application fails during startup with a clear log message instead of returning request-time `500` responses. 

## How to Run Locally

**Prerequisites:** Go 1.24+

```bash
go mod download
go run ./cmd/api
```

The API will be available at:

```text
http://localhost:8080/weather?cep=01001000
```

---

## How to Execute Tests

```bash
go test ./... -v
```

---

## Docker

### Build the image

```bash
docker build -t github.com/Ygohr/fc-weather-cloudrun .
```

### Run the container

```bash
docker run --rm \
  -p 8080:8080 \
  -e PORT=8080 \
  -e WEATHER_API_KEY=your-api-key \
  github.com/Ygohr/fc-weather-cloudrun
```

### Run with Docker Compose

```bash
docker compose up --build
```

With custom parameters:

```bash
# Linux / macOS
PORT=8080 WEATHER_API_KEY=your-api-key docker compose up --build

# Windows PowerShell
$env:PORT=8080; $env:WEATHER_API_KEY="your-api-key"; docker compose up --build

# Windows Cmd
set PORT=8080 && set WEATHER_API_KEY=your-api-key && docker compose up --build
```

---

## API Example

Request:

```http
GET /weather?cep=01001000
```

Success response:

```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

Invalid CEP:

```text
HTTP 422
invalid zipcode
```

CEP not found:

```text
HTTP 404
can not find zipcode
```

Unexpected failure:

```text
HTTP 500
internal server error
```

## Project Structure

```text
.
├── cmd
│   └── api
│       └── main.go
├── infrastructure
│   └── http
│       ├── handler.go
│       ├── handler_test.go
│       ├── model.go
│       ├── router.go
│       └── server.go
├── internal
│   ├── model
│   │   ├── weather.go
│   │   ├── zipcode.go
│   │   └── zipcode_test.go
│   ├── service
│   │   ├── model.go
│   │   ├── viacep.go
│   │   ├── viacep_test.go
│   │   ├── weatherapi.go
│   │   └── weatherapi_test.go
│   ├── usecase
│   │   ├── weather.go
│   │   └── weather_test.go
│   └── util
│       ├── converter.go
│       └── converter_test.go
├── tests
│   ├── helpers.go
│   └── mock_http_client.go
├── Dockerfile
├── docker-compose.yml
├── README.md
├── .env.example
├── .gitignore
├── go.mod
└── go.sum
```
