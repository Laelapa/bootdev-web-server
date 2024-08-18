package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"gitlab.com/demetrius.papas/bootdev-web-server/internal/database"
)

func createUserHandler(w http.ResponseWriter, r *http.Request, db *database.DB) {
	type incomingJSON struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	inc := incomingJSON{}
	err := decoder.Decode(&inc)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	email := inc.Email
	password := inc.Password

	rUsr, err := db.CreateUser(email, password)
	if err != nil {
		if errors.Is(err, database.ErrEmailInUse) {
			errRes(err, w, "Email already in use", 409)
			return
		}
		errRes(err, w, "Something went wrong", 500)
		return
	}

	res, err := json.Marshal(rUsr)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(res)

}
