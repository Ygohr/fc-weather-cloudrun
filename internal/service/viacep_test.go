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

func TestViaCEPClientSuccess(t *testing.T) {
	client := NewViaCEPClient(tests.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "https://viacep.com.br/ws/01001000/json/", req.URL.String())
		return tests.JSONResponse(http.StatusOK, `{"localidade":"Sao Paulo"}`), nil
	}))

	city, err := client.GetCityByZipcode(context.Background(), "01001000")

	require.NoError(t, err)
	assert.Equal(t, "Sao Paulo", city)
}

func TestViaCEPClientNotFound(t *testing.T) {
	client := NewViaCEPClient(tests.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return tests.JSONResponse(http.StatusOK, `{"erro":true}`), nil
	}))

	_, err := client.GetCityByZipcode(context.Background(), "99999999")

	assert.ErrorIs(t, err, model.ErrZipcodeNotFound)
}

func TestViaCEPClientNotFoundWithStringError(t *testing.T) {
	client := NewViaCEPClient(tests.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return tests.JSONResponse(http.StatusOK, `{"erro":"true"}`), nil
	}))

	_, err := client.GetCityByZipcode(context.Background(), "23036321")

	assert.ErrorIs(t, err, model.ErrZipcodeNotFound)
}

func TestViaCEPClientMissingCity(t *testing.T) {
	client := NewViaCEPClient(tests.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return tests.JSONResponse(http.StatusOK, `{}`), nil
	}))

	_, err := client.GetCityByZipcode(context.Background(), "23036321")

	assert.ErrorIs(t, err, model.ErrZipcodeNotFound)
}

func TestViaCEPClientInvalidResponse(t *testing.T) {
	client := NewViaCEPClient(tests.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return tests.JSONResponse(http.StatusOK, `not-json`), nil
	}))

	_, err := client.GetCityByZipcode(context.Background(), "01001000")

	assert.Error(t, err)
}

func TestViaCEPClientExternalFailure(t *testing.T) {
	expectedErr := errors.New("network failure")
	client := NewViaCEPClient(tests.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return nil, expectedErr
	}))

	_, err := client.GetCityByZipcode(context.Background(), "01001000")

	assert.ErrorIs(t, err, expectedErr)
}
