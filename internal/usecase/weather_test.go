package usecase

import (
	"context"
	"errors"
	"testing"

	domain "github.com/Ygohr/fc-weather-cloudrun/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type zipcodeServiceMock struct {
	city string
	err  error
}

func (m zipcodeServiceMock) GetCityByZipcode(ctx context.Context, zipcode string) (string, error) {
	return m.city, m.err
}

type weatherServiceMock struct {
	tempC float64
	err   error
}

func (m weatherServiceMock) GetCurrentTemperature(ctx context.Context, city string) (float64, error) {
	return m.tempC, m.err
}

func TestWeatherUseCaseSuccess(t *testing.T) {
	useCase := NewWeatherUseCase(
		zipcodeServiceMock{city: "Sao Paulo"},
		weatherServiceMock{tempC: 28.5},
	)

	weather, err := useCase.GetWeatherByZipcode(context.Background(), "01001000")

	require.NoError(t, err)
	assert.Equal(t, 28.5, weather.TempC)
	assert.Equal(t, 83.30000000000001, weather.TempF)
	assert.Equal(t, 301.5, weather.TempK)
}

func TestWeatherUseCaseInvalidZipcode(t *testing.T) {
	useCase := NewWeatherUseCase(zipcodeServiceMock{}, weatherServiceMock{})

	_, err := useCase.GetWeatherByZipcode(context.Background(), "invalid")

	assert.ErrorIs(t, err, domain.ErrInvalidZipcode)
}

func TestWeatherUseCaseZipcodeNotFound(t *testing.T) {
	useCase := NewWeatherUseCase(
		zipcodeServiceMock{err: domain.ErrZipcodeNotFound},
		weatherServiceMock{},
	)

	_, err := useCase.GetWeatherByZipcode(context.Background(), "01001000")

	assert.ErrorIs(t, err, domain.ErrZipcodeNotFound)
}

func TestWeatherUseCaseWeatherAPIFailure(t *testing.T) {
	expectedErr := errors.New("weather failed")
	useCase := NewWeatherUseCase(
		zipcodeServiceMock{city: "Sao Paulo"},
		weatherServiceMock{err: expectedErr},
	)

	_, err := useCase.GetWeatherByZipcode(context.Background(), "01001000")

	assert.ErrorIs(t, err, expectedErr)
}
