package main

import (
	"fmt"
	"net/http"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {
	fmt.Println("Server booting...")
	mux := http.NewServeMux()
	mux.Handle("/app/*", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", readinessHandler)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("### SERVER RUNNING ###")
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println("Error on ListenAndServe() call:", err)
	}
}
