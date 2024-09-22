# current-weather

`current-weather` is a simple http service to retrieve a general description of the current weather at the location given in latitude/longitude.  The service is provided on port 8080.

Only a single endpoint is supported which takes lat,long parameters in floating point format.

Examples:
```
$ go run ./cmd/current-weather &

# Palm Springs, CA
$ curl http://localhost:8080?lat=33.822666&long=-116.531418 | jq .
{
  "perception": "Really hot",
  "temperature": 101,
  "shortForecast": "Sunny"
}

# Punxsutawney, PA
$ curl http://localhost:8080?lat=40.945454&long=-78.975175 | jq .
{
  "perception": "Comfortably warm",
  "temperature": 79,
  "shortForecast": "Mostly Sunny"
}
```