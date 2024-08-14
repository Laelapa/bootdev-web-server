package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"gitlab.com/demetrius.papas/bootdev-web-server/internal/database"
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

func fetchChirps(w http.ResponseWriter, r *http.Request, db *database.DB) {
	chirps, err := db.GetChirps()
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	res, err := json.Marshal(chirps)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(res)
}

func validateChirp(w http.ResponseWriter, r *http.Request, db *database.DB) {
	type incomingJSON struct {
		Body string `json:"body"`
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

	c := inc.Body

	for _, v := range []string{"kerfuffle", "sharbert", "fornax"} {
		c = profanityFilter(c, v)
	}

	chrp, err := db.CreateChirp(c)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	res, err := json.Marshal(chrp)
	if err != nil {
		errRes(err, w, "Something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(res)
}

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

func profanityFilter(strToFilter string, badword string) string {
	s := strings.Split(strToFilter, " ")
	for i, str := range s {
		if strings.ToLower(str) == badword {
			s[i] = strings.Replace(strings.ToLower(str), badword, "****", -1)
		}
	}

	return strings.Join(s, " ")

}

func main() {
	fmt.Println("Server booting...")
	db, err := database.NewDB("database.json")
	if err != nil {
		fmt.Println("Error on trying to init the db: ", err)
		return
	}

	var apiCfg apiConfig
	mux := http.NewServeMux()
	fServer := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", apiCfg.middlewareMectricsInc(fServer))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) { validateChirp(w, r, db) })
	mux.HandleFunc("GET /api/chirps", func(w http.ResponseWriter, r *http.Request) { fetchChirps(w, r, db) })
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
