package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"gitlab.com/demetrius.papas/bootdev-web-server/internal/authentication"
	"gitlab.com/demetrius.papas/bootdev-web-server/internal/database"
)

func refreshHandler(w http.ResponseWriter, r *http.Request, db *database.DB) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	hParts := strings.Split(authHeader, " ")
	if len(hParts) != 2 || hParts[0] != "Bearer" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	token := hParts[1]

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
