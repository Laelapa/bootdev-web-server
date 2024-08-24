package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/demetrius.papas/bootdev-web-server/internal/database"
)

func (apiCfg *apiConfig) webhookChirpyRed(w http.ResponseWriter, r *http.Request, db *database.DB) {
	token, err := extractAuthToken(r)
	fmt.Printf("@webhookChirpyRed - token: %v\n", token)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if token != apiCfg.polka {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type incomingJSON struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	inc := incomingJSON{}
	err = decoder.Decode(&inc)
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
