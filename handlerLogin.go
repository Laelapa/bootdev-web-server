package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Laelapa/bootdev-web-server/internal/database"
)

func loginHandler(w http.ResponseWriter, r *http.Request, db *database.DB) {
	type incomingJSON struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	inc := incomingJSON{}
	err := decoder.Decode(&inc)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	user, err := db.LoginUser(inc.Email, inc.Password, inc.ExpiresInSeconds)
	if err != nil {
		if errors.Is(err, database.ErrWrongCredentials) {
			errRes(err, w, "Wrong email or password", 401)
			return
		}
		errRes(err, w, "Something went wrong", 500)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
