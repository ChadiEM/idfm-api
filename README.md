# Île-de-France Mobilités API (RATP, SNCF)

A simple, minimal, and easy-to-understand API to get the timings of Île-de-France Mobilités public transport.

## Build

`go build -o idfm ./cmd/idfm`

## API key

The API key can be retrieved from [here](https://prim.iledefrance-mobilites.fr/fr/mes-jetons-authentification).

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