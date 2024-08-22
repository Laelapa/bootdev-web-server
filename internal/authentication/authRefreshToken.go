package authentication

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const refKeySize = 32 // 32 bytes - 256 bits

// Returns a hex-encoded string fit to be used as a refresh token
func GenerateRefreshToken() (string, error) {
	refTokenBytes := make([]byte, refKeySize)
	_, err := rand.Read(refTokenBytes)
	if err != nil {
		fmt.Printf("Error while generating Refresh Token: %v\n", err)
		return "", err
	}

	refToken := hex.EncodeToString(refTokenBytes)
	return refToken, nil

}
