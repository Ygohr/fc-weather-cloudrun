package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Ygohr/fc-weather-cloudrun/internal/model"
	"github.com/Ygohr/fc-weather-cloudrun/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeatherAPIClientSuccess(t *testing.T) {
	client := NewWeatherAPIClient(tests.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "test-key", req.URL.Query().Get("key"))
		assert.Equal(t, "Sao Paulo", req.URL.Query().Get("q"))
		assert.Equal(t, "no", req.URL.Query().Get("aqi"))
		return tests.JSONResponse(http.StatusOK, `{"current":{"temp_c":28.5}}`), nil
	}), "test-key")

	tempC, err := client.GetCurrentTemperature(context.Background(), "Sao Paulo")

	require.NoError(t, err)
	assert.Equal(t, 28.5, tempC)
}

func TestWeatherAPIClientExternalFailure(t *testing.T) {
	client := NewWeatherAPIClient(tests.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return tests.JSONResponse(http.StatusInternalServerError, `{"error":"failure"}`), nil
	}), "test-key")

	_, err := client.GetCurrentTemperature(context.Background(), "Sao Paulo")

	assert.ErrorIs(t, err, model.ErrWeatherFailure)
}

func TestWeatherAPIClientNetworkFailure(t *testing.T) {
	expectedErr := errors.New("network failure")
	client := NewWeatherAPIClient(tests.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return nil, expectedErr
	}), "test-key")

	_, err := client.GetCurrentTemperature(context.Background(), "Sao Paulo")

	assert.ErrorIs(t, err, expectedErr)
}
