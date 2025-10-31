package auth

import (
	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"net/http"
	"strings"
	"errors"
)

// HashPassword creates an Argon2id hash of the plain-text password.
func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

// CheckPasswordHash compares a plain-text password with a hash.
func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	return match, err
}


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer: "chirpy",
		Subject: userID.String(),
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(tokenSecret), nil

	})

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, jwt.ErrTokenInvalidAudience
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil

}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

    // Split on whitespace, ignore multiple spaces
	parts := strings.Fields(authHeader)
	if len(parts) < 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("authorization header must be in 'Bearer <token>' format")
	}

	token := parts[1]
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}