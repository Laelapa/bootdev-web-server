package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
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
	w.Write([]byte("OK"))
}
