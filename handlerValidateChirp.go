package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"gitlab.com/demetrius.papas/bootdev-web-server/internal/database"
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

	chrp, err := db.CreateChirp(c)
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
