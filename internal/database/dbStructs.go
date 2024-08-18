package database

import (
	"sort"
	"strconv"
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
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

func (u User) GetID() int {
	return u.ID
}

type UserR struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (u UserR) GetID() int {
	return u.ID
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
