# Île-de-France Mobilités API (RATP, SNCF)

A simple, minimal, and easy-to-understand API to get the timings of Île-de-France Mobilités public transport.

## Build

`go build -o idfm ./cmd/idfm`

## API key

The API key can be retrieved from [here](https://prim.iledefrance-mobilites.fr/fr/mes-jetons-authentification).

## Start

To start, `IDFM_API_KEY=<your-api-key> ./idfm`

## Request

`curl http://localhost:8080/api/idfm/timings/bus/42/Versailles%20-%20Chardon%20Lagache?direction=R`

## Sample response

```json
[
  {
    "dest": "Gare Saint-Lazare",
    "time": "2 mn",
    "status": "onTime"
  },
  {
    "dest": "Gare Saint-Lazare",
    "time": "14 mn",
    "status": "onTime"
  }
]
```

## Other examples

### RER A, Auber, all directions

`curl "http://localhost:8080/api/idfm/timings/rail/A/Auber"`

### RER C, Pont du Garigliano, Platform 1

`curl "http://localhost:8080/api/idfm/timings/rail/C/Pont%20du%20Garigliano%20-%20H%C3%B4pital%20Europ%C3%A9en%20G.%20Pompidou?platform=1"`

### Metro 9, Exelmans, direction Aller (Montreuil)

`curl "http://localhost:8080/api/idfm/timings/metro/9/Exelmans?direction=A"`
