package database

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
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
