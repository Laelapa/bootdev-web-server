package main

import (
	"errors"
	"net/http"
	"strings"
)

var ErrWrongHeaderFormat = errors.New("wrong header format")

func extractAuthToken(r *http.Request) (token string, err error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}

	hParts := strings.Split(authHeader, " ")
	if len(hParts) != 2 || hParts[0] != "Bearer" {
		return "", ErrWrongHeaderFormat
	}

	token = hParts[1]
	return token, nil

}
