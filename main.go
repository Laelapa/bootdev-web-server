package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server booting...")
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	srv := &http.Server {
		Addr: ":8080",
		Handler: mux,
	}
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println("Error on ListenAndServe() call:", err)
	}
	fmt.Println("### Server running ###")

}
