# Île-de-France Mobilités API (RATP, SNCF)

A simple, minimal, and easy-to-understand API to get the timings of Île-de-France Mobilités public transport.

## Build

`go build -o idfm ./cmd/idfm`

## Start

To start, `IDFM_API_KEY=<your-api-key> ./idfm`

## Request

`curl localhost:8080/api/idfm/timings/bus/42/Versailles - Chardon Lagache/R`

## Sample response

```json
[
  {
    "dest": "Gare Saint-Lazare",
    "time": "13 mn"
  },
  {
    "dest": "Gare Saint-Lazare",
    "time": "26 mn"
  }
]
```