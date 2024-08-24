package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrWrongHeaderFormat = errors.New("wrong header format")

func extractAuthToken(r *http.Request) (token string, err error) {
	authHeader := r.Header.Get("Authorization")
	fmt.Printf("@extractAuthToken - authHeader: %v\n", authHeader)
	if authHeader == "" {
		return "", nil
	}

	hParts := strings.Split(authHeader, " ")
	if len(hParts) != 2 || hParts[0] != "ApiKey" {
		return "", ErrWrongHeaderFormat
	}

	token = hParts[1]
	return token, nil

}
