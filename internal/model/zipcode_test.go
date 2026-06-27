package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateZipcode(t *testing.T) {
	t.Run("valid CEP", func(t *testing.T) {
		assert.NoError(t, ValidateZipcode("01001000"))
	})

	t.Run("invalid CEP", func(t *testing.T) {
		assert.ErrorIs(t, ValidateZipcode("01001-000"), ErrInvalidZipcode)
		assert.ErrorIs(t, ValidateZipcode("123"), ErrInvalidZipcode)
		assert.ErrorIs(t, ValidateZipcode("abcdefgh"), ErrInvalidZipcode)
	})
}
