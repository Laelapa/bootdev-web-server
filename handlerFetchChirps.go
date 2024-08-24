package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gitlab.com/demetrius.papas/bootdev-web-server/internal/database"
)

func fetchChirps(w http.ResponseWriter, r *http.Request, db *database.DB) {
	authorID, err := strconv.Atoi(r.URL.Query().Get("author_id"))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if authorID > 0 { // Return chirps of specified user
		chirps, err := db.GetChirpsByUser(authorID)
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

		return

	} else if authorID < 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return

	} else { // Return all chirps

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
}
