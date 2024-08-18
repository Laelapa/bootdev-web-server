package main

import (
	"encoding/json"
	"net/http"

	"gitlab.com/demetrius.papas/bootdev-web-server/internal/database"
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

	successfulLogin, err := db.LoginUser(inc.Email, inc.Password)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	if !successfulLogin {
		errRes(err, w, "Wrong credentials", 401)
		return
	}

	user, err := db.GetUserByEmailSanitized(inc.Email)
	if err != nil {
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
