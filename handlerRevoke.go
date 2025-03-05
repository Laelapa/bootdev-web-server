package main

import (
	"errors"
	"net/http"

	"github.com/Laelapa/bootdev-web-server/internal/database"
)

func revokeHandler(w http.ResponseWriter, r *http.Request, db *database.DB) {
	token, err := extractAuthToken(r)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = db.RevokeUserRefreshToken(token)
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
