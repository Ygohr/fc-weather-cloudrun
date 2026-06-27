package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpinfra "github.com/Ygohr/fc-weather-cloudrun/infrastructure/http"
	"github.com/Ygohr/fc-weather-cloudrun/internal/service"
	"github.com/Ygohr/fc-weather-cloudrun/internal/usecase"
)

func main() {
	port := getenv("PORT", "8080")
	apiKey := requiredEnv("WEATHER_API_KEY")

	viaCEPClient := service.NewViaCEPClient(nil)
	weatherAPIClient := service.NewWeatherAPIClient(nil, apiKey)
	weatherUseCase := usecase.NewWeatherUseCase(viaCEPClient, weatherAPIClient)
	handler := httpinfra.NewHandler(weatherUseCase)
	router := httpinfra.NewRouter(handler)
	server := httpinfra.NewServer(port, router)

	go func() {
		log.Printf("server listening on port %s", port)
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func requiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s environment variable is required", key)
	}

	return value
}
