package main

import (
	"net/http"
	"strconv"
	"strings"

	"gitlab.com/demetrius.papas/bootdev-web-server/internal/authentication"
	"gitlab.com/demetrius.papas/bootdev-web-server/internal/database"
)

func chirpDeleter(w http.ResponseWriter, r *http.Request, db *database.DB) {
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

	uID, err := authentication.ValidateJWT(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	chirp, err := db.GetChirp(chirpID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if chirp.AuthorID != uID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err = db.DeleteChirp(chirpID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
