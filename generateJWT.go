package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const oneDay = 84600
const defaultDuration = oneDay

// If time set by user > 24h OR isn't set OR is invalid return DEFAULT_DURATION
// Else return the time set by user
func validateExpiration(expiresInSeconds int) int {
	if expiresInSeconds > defaultDuration || expiresInSeconds <= 0 {
		return defaultDuration
	}

	return expiresInSeconds
}

func generateJWT(userID, expiresInSeconds int, secretKey string) (string, error) {

	expiresInSeconds = validateExpiration(expiresInSeconds)

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(expiresInSeconds) * time.Second)),
		Subject:   strconv.Itoa(userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
