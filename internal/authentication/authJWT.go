package authentication

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
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

func GenerateJWT(userID, expiresInSeconds int) (string, error) {

	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")

	expiresInSeconds = validateExpiration(expiresInSeconds)

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(expiresInSeconds) * time.Second)),
		Subject:   strconv.Itoa(userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func ValidateJWT(tokenStr string) (int, error) {
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {

		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("wrong signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}

	uIDstr, _ := token.Claims.GetSubject()
	uID, _ := strconv.Atoi(uIDstr)

	return uID, nil
}
