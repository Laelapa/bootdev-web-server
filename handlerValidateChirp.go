package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Laelapa/bootdev-web-server/internal/authentication"
	"github.com/Laelapa/bootdev-web-server/internal/database"
)

func validateChirp(w http.ResponseWriter, r *http.Request, db *database.DB) {
	type incomingJSON struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	inc := incomingJSON{}
	err := decoder.Decode(&inc)

	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	if len(inc.Body) > 140 {
		errRes(errors.New("chirp is too long"), w, "Chirp is too long", 400)
		return
	}

	c := inc.Body

	for _, v := range []string{"kerfuffle", "sharbert", "fornax"} {
		c = profanityFilter(c, v)
	}

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

	chrp, err := db.CreateChirp(c, uID)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	res, err := json.Marshal(chrp)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(res)
}
