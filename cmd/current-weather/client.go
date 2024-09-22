package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// URL for the National Weather Service (NWS) API
	weatherBaseURL = "https://api.weather.gov"
)

// WeatherClient provides limited access to the NWS API
type WeatherClient struct {
	baseurl string
}

// NewWeatherClient creates the client initialized with the NWS API URL
func NewWeatherClient() WeatherClient {
	return WeatherClient{
		baseurl: weatherBaseURL,
	}
}

// ForecastResponse provides the structure to decode data returned for the NWS forecast
type ForecastResponse struct {
	Properties struct {
		Periods []struct {
			Number        int    `json:"number"`
			Name          string `json:"name"`
			Temperature   int    `json:"temperature"`
			ShortForecast string `json:"shortForecast"`
		}
	} `json:"properties"`
}

// GetForecast gets the forecast for the given location (latitude/longitude)
func (wc WeatherClient) GetForecast(lat, long float64) (*ForecastResponse, error) {
	// convert lat/long to gridpoint
	url, err := wc.getForecastURL(lat, long)
	if err != nil {
		return nil, fmt.Errorf("getting gridpoint: %w", err)
	}

	// get forecast from url provided in response
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("getting forecast from %s: %w", url, err)
	}

	fc := ForecastResponse{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&fc); err != nil {
		return nil, fmt.Errorf("decoding forecast json: %w", err)
	}

	return &fc, nil
}

// PointsResponse provides the structure to decode data returned for the NWS location
type PointsResponse struct {
	Properties struct {
		// URL to retrieve forecast
		Forecast string
	} `json:"properties"`
}

// getForecastURL retrieves the URL to get the forecase for the given latitude/longitude
func (wc WeatherClient) getForecastURL(lat, long float64) (string, error) {
	url := fmt.Sprintf("%s/points/%0.3f,%0.3f", wc.baseurl, lat, long)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("calling %s: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("calling %s: %s", url, resp.Status)
	}

	pr := PointsResponse{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&pr); err != nil {
		return "", fmt.Errorf("decoding points response: %w", err)
	}

	return pr.Properties.Forecast, nil
}
