package auth

import (
	"github.com/alexedwards/argon2id"
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