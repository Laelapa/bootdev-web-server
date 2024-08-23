package main

import (
	"encoding/json"
	"net/http"

	"gitlab.com/demetrius.papas/bootdev-web-server/internal/database"
)

func webhookChirpyRed(w http.ResponseWriter, r *http.Request, db *database.DB) {
	type incomingJSON struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	inc := incomingJSON{}
	err := decoder.Decode(&inc)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	if inc.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	isUpgraded, err := db.UserP2W(inc.Data.UserID)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	if isUpgraded {
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
