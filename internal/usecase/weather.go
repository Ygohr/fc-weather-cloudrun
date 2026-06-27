package usecase

import (
	"context"

	domain "github.com/Ygohr/fc-weather-cloudrun/internal/model"
	"github.com/Ygohr/fc-weather-cloudrun/internal/util"
)

type ZipcodeService interface {
	GetCityByZipcode(ctx context.Context, zipcode string) (string, error)
}

type WeatherService interface {
	GetCurrentTemperature(ctx context.Context, city string) (float64, error)
}

type WeatherUseCase struct {
	zipcodeService ZipcodeService
	weatherService WeatherService
}

func NewWeatherUseCase(zipcodeService ZipcodeService, weatherService WeatherService) *WeatherUseCase {
	return &WeatherUseCase{
		zipcodeService: zipcodeService,
		weatherService: weatherService,
	}
}

func (u *WeatherUseCase) GetWeatherByZipcode(ctx context.Context, zipcode string) (domain.Weather, error) {
	if err := domain.ValidateZipcode(zipcode); err != nil {
		return domain.Weather{}, err
	}

	city, err := u.zipcodeService.GetCityByZipcode(ctx, zipcode)
	if err != nil {
		return domain.Weather{}, err
	}

	tempC, err := u.weatherService.GetCurrentTemperature(ctx, city)
	if err != nil {
		return domain.Weather{}, err
	}

	return domain.Weather{
		TempC: tempC,
		TempF: util.CelsiusToFahrenheit(tempC),
		TempK: util.CelsiusToKelvin(tempC),
	}, nil
}
