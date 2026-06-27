package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Ygohr/fc-weather-cloudrun/internal/model"
)

type ViaCEPClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewViaCEPClient(httpClient *http.Client) *ViaCEPClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 5 * time.Second}
	}

	return &ViaCEPClient{
		httpClient: httpClient,
		baseURL:    "https://viacep.com.br/ws",
	}
}

func NewViaCEPClientWithBaseURL(httpClient *http.Client, baseURL string) *ViaCEPClient {
	client := NewViaCEPClient(httpClient)
	client.baseURL = baseURL
	return client
}

func (c *ViaCEPClient) GetCityByZipcode(ctx context.Context, zipcode string) (string, error) {
	requestURL := fmt.Sprintf("%s/%s/json/", c.baseURL, url.PathEscape(zipcode))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", model.ErrZipcodeNotFound
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("viacep returned status %d", resp.StatusCode)
	}

	var data viaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if data.IsNotFound() {
		return "", model.ErrZipcodeNotFound
	}

	if data.Localidade == "" {
		return "", model.ErrZipcodeNotFound
	}

	return data.Localidade, nil
}
