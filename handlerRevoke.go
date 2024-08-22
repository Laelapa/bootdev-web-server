package main

import (
	"errors"
	"net/http"
	"strings"

	"gitlab.com/demetrius.papas/bootdev-web-server/internal/database"
)

func revokeHandler(w http.ResponseWriter, r *http.Request, db *database.DB) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	hParts := strings.Split(authHeader, " ")
	if len(hParts) != 2 || hParts[0] != "Bearer" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	token := hParts[1]

	err := db.RevokeUserRefreshToken(token)
	if err != nil {
		if errors.Is(err, database.ErrInvalidRefToken) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		} else {
			errRes(err, w, "Something went wrong", 500)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
