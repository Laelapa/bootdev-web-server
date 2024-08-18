package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func errRes(err error, w http.ResponseWriter, report string, errCode int) {
	type errorJSON struct {
		ErrorBody string `json:"error"`
	}

	fmt.Printf("%s", err)
	errJSON := errorJSON{
		ErrorBody: report,
	}
	res, err := json.Marshal(errJSON)
	if err != nil {
		fmt.Println("Error trying to send an error response: ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errCode)
	w.Write(res)
}
