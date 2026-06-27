package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	domain "github.com/Ygohr/fc-weather-cloudrun/internal/model"
)

type WeatherAPIClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}



func NewWeatherAPIClient(httpClient *http.Client, apiKey string) *WeatherAPIClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 5 * time.Second}
	}

	return &WeatherAPIClient{
		httpClient: httpClient,
		apiKey:     apiKey,
		baseURL:    "https://api.weatherapi.com/v1/current.json",
	}
}

func NewWeatherAPIClientWithBaseURL(httpClient *http.Client, apiKey string, baseURL string) *WeatherAPIClient {
	client := NewWeatherAPIClient(httpClient, apiKey)
	client.baseURL = baseURL
	return client
}

func (c *WeatherAPIClient) GetCurrentTemperature(ctx context.Context, city string) (float64, error) {
	if c.apiKey == "" {
		return 0, errors.New("weather api key is required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
	if err != nil {
		return 0, err
	}

	query := req.URL.Query()
	query.Set("key", c.apiKey)
	query.Set("q", city)
	query.Set("aqi", "no")
	req.URL.RawQuery = query.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return 0, fmt.Errorf("%w: status %d", domain.ErrWeatherFailure, resp.StatusCode)
	}

	var data weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.Current.TempC, nil
}
