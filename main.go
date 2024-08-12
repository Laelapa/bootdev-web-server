package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMectricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {
	fmt.Println("Server booting...")
	var apiCfg apiConfig
	mux := http.NewServeMux()
	fServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", apiCfg.middlewareMectricsInc(fServer))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /api/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetHandler)
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
