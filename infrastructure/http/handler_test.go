package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ygohr/fc-weather-cloudrun/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type weatherUseCaseMock struct {
	weather model.Weather
	err     error
}

func (m weatherUseCaseMock) GetWeatherByZipcode(ctx context.Context, zipcode string) (model.Weather, error) {
	return m.weather, m.err
}

func TestHandlerWeather200(t *testing.T) {
	handler := NewHandler(weatherUseCaseMock{weather: model.Weather{
		TempC: 28.5,
		TempF: 83.3,
		TempK: 301.5,
	}})
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=01001000", nil)
	rec := httptest.NewRecorder()

	handler.Weather(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"temp_C":28.5,"temp_F":83.3,"temp_K":301.5}`, rec.Body.String())
}

func TestHandlerWeather404(t *testing.T) {
	handler := NewHandler(weatherUseCaseMock{err: model.ErrZipcodeNotFound})
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=99999999", nil)
	rec := httptest.NewRecorder()

	handler.Weather(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "can not find zipcode", strings.TrimSpace(rec.Body.String()))
}

func TestHandlerWeather422(t *testing.T) {
	handler := NewHandler(weatherUseCaseMock{err: model.ErrInvalidZipcode})
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=invalid", nil)
	rec := httptest.NewRecorder()

	handler.Weather(rec, req)

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Equal(t, "invalid zipcode", strings.TrimSpace(rec.Body.String()))
}

func TestHandlerWeather500(t *testing.T) {
	handler := NewHandler(weatherUseCaseMock{err: errors.New("unexpected")})
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=01001000", nil)
	rec := httptest.NewRecorder()

	handler.Weather(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "internal server error", strings.TrimSpace(rec.Body.String()))
}
