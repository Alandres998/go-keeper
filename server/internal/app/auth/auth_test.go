package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	userID := 1
	ctx := context.Background()

	token, err := GenerateToken(ctx, userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	id, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, id)
}
