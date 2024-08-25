package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strconv"

	"gitlab.com/demetrius.papas/bootdev-web-server/internal/database"
)

func fetchChirps(w http.ResponseWriter, r *http.Request, db *database.DB) {
	var err error

	authorID := r.URL.Query().Get("author_id")
	iauthorID := 0

	if authorID != "" {
		iauthorID, err = strconv.Atoi(authorID)
		if err != nil {
			http.Error(w, "Bad Request wtf bro", http.StatusBadRequest)
			return
		}
	}

	sortQ := r.URL.Query().Get("sort")

	if iauthorID > 0 { // Return chirps of specified user
		chirps, err := db.GetChirpsByUser(iauthorID)
		if err != nil {
			errRes(err, w, "Something went wrong", 500)
			return
		}

		if sortQ == "desc" {
			slices.Reverse(chirps)
		}

		res, err := json.Marshal(chirps)
		if err != nil {
			errRes(err, w, "Something went wrong", 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)

		return

	} else if iauthorID < 0 {
		http.Error(w, "Bad Request eyy 0", http.StatusBadRequest)
		return

	} else { // Return all chirps

		chirps, err := db.GetChirps()
		if err != nil {
			errRes(err, w, "Something went wrong", 500)
			return
		}

		if sortQ == "desc" {
			slices.Reverse(chirps)
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
