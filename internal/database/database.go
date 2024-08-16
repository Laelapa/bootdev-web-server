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

type DBEntry interface {
	GetID() int
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (c Chirp) GetID() int {
	return c.ID
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (u User) GetID() int {
	return u.ID
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type DBStructureStr struct {
	Chirps map[string]Chirp `json:"chirps"`
	Users  map[string]User  `json:"users"`
}

// Convert a DBStructure to a DBStructureStr, which contains string keys instead of int
func (db *DBStructure) mapItoa() *DBStructureStr {
	itoa := DBStructureStr{
		Chirps: make(map[string]Chirp),
		Users:  make(map[string]User),
	}

	for i, v := range db.Chirps {
		itoa.Chirps[strconv.Itoa(i)] = v
	}

	for i, v := range db.Users {
		itoa.Users[strconv.Itoa(i)] = v
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
	mDB.Users = make(map[int]User)

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

// Returns a single chirp with matching ID. If none found returns zero value and nil error
func (db *DB) GetChirp(id int) (Chirp, error) {
	chirps, err := db.GetChirps()
	if err != nil {
		return Chirp{}, fmt.Errorf("%w", err)
	}

	for _, v := range chirps {
		if v.ID == id {
			return v, nil
		}
	}

	return Chirp{}, nil
}

// Returns a single user with matching ID. If none found returns zero value and nil error
func (db *DB) GetUser(id int) (User, error) {
	users, err := db.GetUsers()
	if err != nil {
		return User{}, fmt.Errorf("%w", err)
	}

	for _, v := range users {
		if v.ID == id {
			return v, nil
		}
	}

	return User{}, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	DBs, err := db.loadDB()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	chrp := toSlice(DBs.Chirps)
	return chrp, nil
}

func (db *DB) GetUsers() ([]User, error) {
	DBs, err := db.loadDB()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	usr := toSlice(DBs.Users)
	return usr, nil
}

// Returns a, sorted by id, slice of all the DBEntries in the structure.
// Leaves space in the underlying array for one more chirp
func toSlice[T DBEntry](entries map[int]T) []T {
	sliceOfEntries := make([]T, 0, len(entries)+1)
	for _, v := range entries {
		// append instead of e[i] = v to avoid panic due to line 152
		sliceOfEntries = append(sliceOfEntries, v)
	}

	sort.Slice(sliceOfEntries, func(a, b int) bool { return sliceOfEntries[a].GetID() < sliceOfEntries[b].GetID() })

	return sliceOfEntries
}

// Creates a new chirp in the database and also returns it
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, fmt.Errorf("%w", err)
	}

	chirps := toSlice(dbStructure.Chirps)

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

// Creates a new user in the database and also returns it
func (db *DB) CreateUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, fmt.Errorf("%w", err)
	}

	users := toSlice(dbStructure.Users)

	usr := User{
		ID:    len(users) + 1,
		Email: email,
	}
	dbStructure.Users[usr.ID] = usr

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, fmt.Errorf("%w", err)
	}

	return usr, nil
}
