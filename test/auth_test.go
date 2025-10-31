package test

import (
	"testing"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"chirpy/internal/auth"

)

const testSecret = "super-secret-jwt-key-for-testing"

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	expiresIn := 1 * time.Hour

	// Create JWT
	tokenStr, err := auth.MakeJWT(userID, testSecret, expiresIn)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	// Validate JWT
	parsedID, err := auth.ValidateJWT(tokenStr, testSecret)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedID)
}

func TestValidateJWT_Expired(t *testing.T) {
	userID := uuid.New()
	expiresIn := -1 * time.Second // already expired

	tokenStr, err := auth.MakeJWT(userID, testSecret, expiresIn)
	assert.NoError(t, err)

	parsedID, err := auth.ValidateJWT(tokenStr, testSecret)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
	assert.Equal(t, uuid.Nil, parsedID)
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	expiresIn := 1 * time.Hour

	tokenStr, err := auth.MakeJWT(userID, testSecret, expiresIn)
	assert.NoError(t, err)

	// Try with wrong secret
	parsedID, err := auth.ValidateJWT(tokenStr, "wrong-secret")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature is invalid")
	assert.Equal(t, uuid.Nil, parsedID)
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	parsedID, err := auth.ValidateJWT("not.a.real.token", testSecret)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, parsedID)
}

func TestValidateJWT_InvalidSubject(t *testing.T) {
	// Manually craft a token with invalid subject
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   "not-a-uuid",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(testSecret))

	parsedID, err := auth.ValidateJWT(signed, testSecret)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID")
	assert.Equal(t, uuid.Nil, parsedID)
}