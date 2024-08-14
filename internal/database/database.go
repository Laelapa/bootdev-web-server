package database

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"sync"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type DBStructureStr struct {
	Chirps map[string]Chirp `json:"chirps"`
}

// Convert a DBStructure to a DBStructureStr, which contains string keys instead of int
func (db *DBStructure) mapItoa() *DBStructureStr {
	itoa := DBStructureStr{}

	for i, v := range db.Chirps {
		itoa.Chirps[strconv.Itoa(i)] = v
	}

	return &itoa
}

func NewDB(path string) (*DB, error) {
	d := &DB{}
	d.path = path

	err := d.ensureDB()
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	_, err := os.Create(db.path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbFile, err := os.Open(db.path)
	if err != nil {
		return DBStructure{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer dbFile.Close()

	dbFileBytes, err := io.ReadAll(dbFile)
	if err != nil {
		return DBStructure{}, fmt.Errorf("failed while reading file: %w", err)
	}

	var mDB DBStructure
	json.Unmarshal(dbFileBytes, &mDB)
	return mDB, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	DBSs := dbStructure.mapItoa()

	jsonData, err := json.MarshalIndent(DBSs, "", " ")
	if err != nil {
		return fmt.Errorf("failed while marshalling json: %w", err)
	}

	db.mux.RLock()
	dbFile, err := os.OpenFile(db.path, os.O_APPEND|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open file for write: %w", err)
	}
	defer db.mux.RUnlock()
	defer dbFile.Close()

	_, err = dbFile.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed while writing json data to db file: %w", err)
	}

	fmt.Println("JSON data saved to disk")
	return nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	DBs, err := db.loadDB()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	chrp := DBs.chirpsToSlice()
	return chrp, nil
}

// Returns a, sorted by id, slice of all the chirps in the structure. Leaves space in
// the underlying array for one more chirp
func (dbStructure *DBStructure) chirpsToSlice() []Chirp {
	chrp := make([]Chirp, len(dbStructure.Chirps), len(dbStructure.Chirps)+1)
	for i, v := range dbStructure.Chirps {
		chrp[i] = v
	}

	sort.Slice(chrp, func(a, b int) bool { return chrp[a].ID < chrp[b].ID })

	return chrp
}

// Creates a new chirp in the database and also returns it
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, fmt.Errorf("%w", err)
	}

	chirps := dbStructure.chirpsToSlice()

	var chrp Chirp
	chrp.Body = body
	chrp.ID = len(chirps)
	dbStructure.Chirps[chrp.ID] = chrp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, fmt.Errorf("%w", err)
	}

	return chrp, nil
}
