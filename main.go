package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type apiConfig struct {
	FileserverHits int
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMectricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	tmpl, err := template.New("metrics").Parse(`
		<html>

		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited {{ .FileserverHits }} times!</p>
		</body>

		</html>
	`)
	if err != nil {
		fmt.Println("Error while parsing HTML:", err)
		return
	}
	err = tmpl.Execute(w, cfg)
	if err != nil {
		fmt.Println("Error while executing the template:", err)
		return
	}
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type incomingJSON struct {
		Body string `json:"body"`
	}

	type responseJSON struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	inc := incomingJSON{}
	err := decoder.Decode(&inc)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}
	if len(inc.Body) > 140 {
		errRes(errors.New("chirp is too long"), w, "Chirp is too long", 400)
		return
	}
	resJSON := responseJSON {
		Valid: true,
	}
	res, err := json.Marshal(resJSON)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(res)
}

func errRes(err error, w http.ResponseWriter, report string, errCode int) {
	type errorJSON struct {
		ErrorBody string `json:"error"`
	}

	log.Printf("%s", err)
	errJSON := errorJSON {
		ErrorBody: report,
	}
	res, err := json.Marshal(errJSON)
	if err != nil {
		log.Printf("Error trying to send an error response: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errCode)
	w.Write(res)
}

func main() {
	fmt.Println("Server booting...")
	var apiCfg apiConfig
	mux := http.NewServeMux()
	fServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", apiCfg.middlewareMectricsInc(fServer))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("### SERVER RUNNING ###")
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println("Error on ListenAndServe() call:", err)
		return
	}
}
