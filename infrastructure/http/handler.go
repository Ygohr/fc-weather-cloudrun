package http

import (
	"context"
	"encoding/json"
	"errors"
	stdhttp "net/http"

	domain "github.com/Ygohr/fc-weather-cloudrun/internal/model"
)

type WeatherUseCase interface {
	GetWeatherByZipcode(ctx context.Context, zipcode string) (domain.Weather, error)
}

type Handler struct {
	weatherUseCase WeatherUseCase
}

func NewHandler(weatherUseCase WeatherUseCase) *Handler {
	return &Handler{weatherUseCase: weatherUseCase}
}

func (h *Handler) Weather(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if r.Method != stdhttp.MethodGet {
		w.WriteHeader(stdhttp.StatusMethodNotAllowed)
		return
	}

	zipcode := r.URL.Query().Get("cep")
	weather, err := h.weatherUseCase.GetWeatherByZipcode(r.Context(), zipcode)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stdhttp.StatusOK)
	_ = json.NewEncoder(w).Encode(weatherResponse{
		TempC: weather.TempC,
		TempF: weather.TempF,
		TempK: weather.TempK,
	})
}

func (h *Handler) handleError(w stdhttp.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidZipcode):
		stdhttp.Error(w, "invalid zipcode", stdhttp.StatusUnprocessableEntity)
	case errors.Is(err, domain.ErrZipcodeNotFound):
		stdhttp.Error(w, "can not find zipcode", stdhttp.StatusNotFound)
	default:
		stdhttp.Error(w, "internal server error", stdhttp.StatusInternalServerError)
	}
}
