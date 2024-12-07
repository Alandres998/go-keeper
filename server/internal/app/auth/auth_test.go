package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	userID := 1
	token, err := GenerateToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	id, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, id)
}
