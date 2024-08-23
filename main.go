package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gitlab.com/demetrius.papas/bootdev-web-server/internal/database"
)

type apiConfig struct {
	FileserverHits int
	jwtSecret      string
}

func main() {
	fmt.Println("Server booting...")

	godotenv.Load()

	db, err := database.NewDB("database.json")
	if err != nil {
		fmt.Println("Error on trying to init the db: ", err)
		return
	}

	var apiCfg apiConfig
	apiCfg.jwtSecret = os.Getenv("JWT_SECRET")

	mux := http.NewServeMux()
	fServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", apiCfg.middlewareMectricsInc(fServer))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) { createUserHandler(w, r, db) })
	mux.HandleFunc("PUT /api/users", func(w http.ResponseWriter, r *http.Request) { userUpdateHandler(w, r, db) })
	mux.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) { loginHandler(w, r, db) })
	mux.HandleFunc("POST /api/refresh", func(w http.ResponseWriter, r *http.Request) { refreshHandler(w, r, db) })
	mux.HandleFunc("POST /api/revoke", func(w http.ResponseWriter, r *http.Request) { revokeHandler(w, r, db) })
	mux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) { validateChirp(w, r, db) })
	mux.HandleFunc("GET /api/chirps", func(w http.ResponseWriter, r *http.Request) { fetchChirps(w, r, db) })
	mux.HandleFunc("GET /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) { chirpGetter(w, r, db) })
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) { chirpDeleter(w, r, db) })
	// Webhooks
	mux.HandleFunc("POST /api/polka/webhooks", func(w http.ResponseWriter, r *http.Request) { webhookChirpyRed(w, r, db) })
	
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("### SERVER RUNNING ###")
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println("Error on ListenAndServe() call:", err)
		return
	}
}
