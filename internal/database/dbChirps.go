package database

import "fmt"

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

func (db *DB) GetChirps() ([]Chirp, error) {
	DBs, err := db.loadDB()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	chrp := toSlice(DBs.Chirps)
	return chrp, nil
}

func (db *DB) GetChirpsByUser(uID int) ([]Chirp, error) {
	DBs, err := db.loadDB()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	chirpsByUser := make([]Chirp, 0, len(DBs.Chirps))
	for _, v := range DBs.Chirps {
		if v.AuthorID == uID {
			chirpsByUser = append(chirpsByUser, v)
		}
	}

	return chirpsByUser, nil
}

// Creates a new chirp in the database and also returns it
func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, fmt.Errorf("%w", err)
	}

	chirps := toSlice(dbStructure.Chirps)

	chrp := Chirp{
		ID:   len(chirps) + 1,
		Body: body,
		AuthorID: authorID,
	}
	dbStructure.Chirps[chrp.ID] = chrp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, fmt.Errorf("%w", err)
	}

	return chrp, nil
}

func (db *DB) DeleteChirp(chirpID int) error {
	dbData, err := db.loadDB()
	if err != nil {
		return err
	}

	wIterator := 0
	var wData DBStructure
	wData.Users = dbData.Users

	for _, v := range dbData.Chirps {
		if v.ID == chirpID {
			continue
		}

		wData.Chirps[wIterator] = v
		wIterator++
	}

	err = db.writeDB(wData)
	if err != nil {
		return err
	}

	return nil
}