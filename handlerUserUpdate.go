package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"gitlab.com/demetrius.papas/bootdev-web-server/internal/authentication"
	"gitlab.com/demetrius.papas/bootdev-web-server/internal/database"
)

func userUpdateHandler(w http.ResponseWriter, r *http.Request, db *database.DB) {
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

	type incomingJSON struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	inc := incomingJSON{}
	err = decoder.Decode(&inc)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	usr, err := db.UpdateUserCredentials(uID, inc.Email, inc.Password)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	res, err := json.Marshal(usr)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
