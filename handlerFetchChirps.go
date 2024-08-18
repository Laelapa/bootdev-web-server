package main

import (
	"encoding/json"
	"net/http"

	"gitlab.com/demetrius.papas/bootdev-web-server/internal/database"
)

func fetchChirps(w http.ResponseWriter, _ *http.Request, db *database.DB) {
	chirps, err := db.GetChirps()
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	res, err := json.Marshal(chirps)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
