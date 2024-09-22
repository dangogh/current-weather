package main

// current-weather provides a simple http service returning a general description
// of the current conditions at the given location.  See README.md for details.

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

const (
	latitudeKeyword  = "lat"
	longitudeKeyword = "long"
	weatherBaseURL   = "https://api.weather.gov"
)

// WeatherResponse is the data returned from this service indicating the
// forecast for the given location
type WeatherResponse struct {
	// general description of temperature range
	Perception string `json:"perception"`
	// temperature in Fahrenheit
	Temperature int `json:"temperature"`
	// short description of forecast for the next 12 hours
	ShortForecast string `json:"shortForecast"`
}

func main() {
	log.Fatal(
		http.ListenAndServe(":8080",
			http.HandlerFunc(forecastHandler()),
		),
	)
}

// forecastHandler handles requests to the service using
func forecastHandler() func(w http.ResponseWriter, r *http.Request) {
	client := NewWeatherClient()
	return func(w http.ResponseWriter, r *http.Request) {
		lat, err := strconv.ParseFloat(r.URL.Query().Get(latitudeKeyword), 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid latitude " + r.URL.Query().Get(latitudeKeyword)))
			return
		}

		long, err := strconv.ParseFloat(r.URL.Query().Get(longitudeKeyword), 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid longitude " + r.URL.Query().Get(longitudeKeyword)))
			return
		}

		forecast, err := client.GetForecast(lat, long)
		if err != nil {
			// here I would try to classify the error to provide a more accurate response code
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("getting forecast: " + err.Error()))
			return
		}

		if len(forecast.Properties.Periods) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("no forecast data"))
			return
		}

		period := forecast.Properties.Periods[0]
		wr := WeatherResponse{
			Perception:    temperaturePerception(period.Temperature),
			Temperature:   period.Temperature,
			ShortForecast: period.ShortForecast,
		}

		enc := json.NewEncoder(w)
		if err := enc.Encode(wr); err != nil {
			w.Write([]byte("error encoding to json: " + err.Error()))
		}
	}
}

// temperaturePerception returns a string describing the temperature range
// in familiar terms
func temperaturePerception(t int) string {
	switch {
	case t > 110:
		return "Dangerously hot"
	case t > 90:
		return "Really hot"
	case t > 80:
		return "Hot"
	case t > 70:
		return "Comfortably warm"
	case t > 60:
		return "Warm"
	case t > 50:
		return "Cool"
	case t > 40:
		return "Chilly"
	case t > 30:
		return "Cold"
	case t > 10:
		return "Really cold"
	case t > 0:
		return "Bone-chilling cold"
	default:
		return "Dangerously cold"
	}
}
