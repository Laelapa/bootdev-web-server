package main

import (
	"net/http"
	"strconv"

	"gitlab.com/demetrius.papas/bootdev-web-server/internal/authentication"
	"gitlab.com/demetrius.papas/bootdev-web-server/internal/database"
)

func chirpDeleter(w http.ResponseWriter, r *http.Request, db *database.DB) {
	token, err := extractAuthToken(r)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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
