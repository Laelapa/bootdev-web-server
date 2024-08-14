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
	itoa := DBStructureStr{
		Chirps: make(map[string]Chirp),
	}

	for i, v := range db.Chirps {
		itoa.Chirps[strconv.Itoa(i)] = v
	}

	return &itoa
}

func NewDB(path string) (*DB, error) {
	d := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

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
	fmt.Println("dbFileBytes:", string(dbFileBytes))

	var mDB DBStructure

	mDB.Chirps = make(map[int]Chirp)

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

// Returns a, sorted by id, slice of all the chirps in the structure.
// Leaves space inthe underlying array for one more chirp
func (dbStructure *DBStructure) chirpsToSlice() []Chirp {
	chrp := make([]Chirp, 0, len(dbStructure.Chirps)+1)
	for _, v := range dbStructure.Chirps {
		// append instead of chrp[i] = v to avoid panic due to line 152
		chrp = append(chrp, v) 
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

	chrp := Chirp{
		ID:   len(chirps) + 1,
		Body: body,
	}
	dbStructure.Chirps[chrp.ID] = chrp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, fmt.Errorf("%w", err)
	}

	return chrp, nil
}
