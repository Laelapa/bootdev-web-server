package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"gitlab.com/demetrius.papas/bootdev-web-server/internal/authentication"
	"gitlab.com/demetrius.papas/bootdev-web-server/internal/database"
)

func refreshHandler(w http.ResponseWriter, r *http.Request, db *database.DB) {
	token, err := extractAuthToken(r)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenUserID, err := db.CheckUserRefreshToken(token)
	if err != nil {
		if errors.Is(err, database.ErrInvalidRefToken) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		} else {
			errRes(err, w, "Something went wrong", 500)
			return
		}
	}

	freshJWT, err := authentication.GenerateJWT(tokenUserID, 60)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	resStruct := struct {
		Token string `json:"token"`
	}{freshJWT}

	res, err := json.Marshal(resStruct)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
