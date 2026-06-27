package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCelsiusToFahrenheit(t *testing.T) {
	assert.InDelta(t, 83.3, CelsiusToFahrenheit(28.5), 0.0001)
}

func TestCelsiusToKelvin(t *testing.T) {
	assert.Equal(t, 301.5, CelsiusToKelvin(28.5))
}
