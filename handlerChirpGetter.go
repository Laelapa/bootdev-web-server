package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Laelapa/bootdev-web-server/internal/database"
)

func chirpGetter(w http.ResponseWriter, r *http.Request, db *database.DB) {
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		http.NotFound(w, r)
		return
	}

	chirpIDi, err := strconv.Atoi(chirpID)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	chrp, err := db.GetChirp(chirpIDi)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}
	if chrp.ID == 0 {
		http.NotFound(w, r)
		return
	}

	res, err := json.Marshal(chrp)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
