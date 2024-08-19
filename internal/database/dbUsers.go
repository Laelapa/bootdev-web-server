package database

import (
	"errors"
	"fmt"

	"gitlab.com/demetrius.papas/bootdev-web-server/internal/authentication"

	"golang.org/x/crypto/bcrypt"
)

var ErrEmailInUse = errors.New("email already in use")
var ErrWrongCredentials = errors.New("wrong username or password")

// Returns a single user with matching ID.
// If none found returns zero value and nil error.
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

// Returns a single user with matching email.
// If none found returns zero value and nil error.
func (db *DB) GetUserByEmail(email string) (User, error) {
	users, err := db.GetUsers()
	if err != nil {
		return User{}, fmt.Errorf("%w", err)
	}

	for _, v := range users {
		if v.Email == email {
			return v, nil
		}
	}

	return User{}, nil
}

func (db *DB) GetUserByEmailSanitized(email string) (UserR, error) {
	usr, err := db.GetUserByEmail(email)
	if err != nil {
		return UserR{}, fmt.Errorf("%w", err)
	}

	var usrResp UserR

	usrResp.ID = usr.ID
	usrResp.Email = usr.Email

	return usrResp, nil
}

// Returns `true` if the email already exists in the database.
func (db *DB) CheckEmailTaken(email string) (isTaken bool, err error) {
	usr, err := db.GetUserByEmail(email)
	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	if usr.ID == 0 && usr.Email == "" {
		return false, nil
	}

	return true, nil
}

func (db *DB) LoginUser(email string, pwd string, expiresInSeconds int) (UserR, error) {
	emailExists, err := db.CheckEmailTaken(email)
	if err != nil {
		return UserR{}, fmt.Errorf("%w", err)
	}

	if !emailExists {
		return UserR{}, ErrWrongCredentials
	}

	user, err := db.GetUserByEmail(email)
	if err != nil {
		return UserR{}, fmt.Errorf("%w", err)
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(pwd))
	if err != nil {
		return UserR{}, ErrWrongCredentials
	}

	userSan, err := db.GetUserByEmailSanitized(user.Email)
	if err != nil {
		return UserR{}, fmt.Errorf("%w", err)
	}

	token, err := authentication.GenerateJWT(user.ID, expiresInSeconds)
	if err != nil {
		return UserR{}, fmt.Errorf("%w", err)
	}

	userSan.Token = token

	return userSan, nil
}

func (db *DB) GetUsers() ([]User, error) {
	DBs, err := db.loadDB()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	usr := toSlice(DBs.Users)
	return usr, nil
}

func (db *DB) UpdateUser(uID int, uEmail string, uPwd string) (User, error) {
	users, err := db.GetUsers()
	if err != nil {
		return User{}, fmt.Errorf("%w", err)
	}

	 DBs, err := db.loadDB()
	 if err != nil {
		return User{}, fmt.Errorf("%w", err)
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(uPwd), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("%w", err)
	}

	for i, v := range users {
		if v.ID == uID {
			users[i].Email = uEmail
			users[i].Password = pwdHash
		}

		DBs.Users[i] = users[i]
	}

	err = db.writeDB(DBs)
	if err != nil {
		return User{}, fmt.Errorf("%w", err)
	}

	return User{ ID: uID, Email: uEmail, }, nil
} 

// Creates a new user in the database and also returns it, without including the pwd
func (db *DB) CreateUser(email string, pwd string) (UserR, error) {
	emailTaken, err := db.CheckEmailTaken(email)
	if err != nil {
		return UserR{}, fmt.Errorf("%w", err)
	}

	if emailTaken {
		return UserR{}, ErrEmailInUse
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return UserR{}, fmt.Errorf("%w", err)
	}

	users := toSlice(dbStructure.Users)
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return UserR{}, fmt.Errorf("%w", err)
	}

	usr := User{
		ID:       len(users) + 1,
		Email:    email,
		Password: pwdHash,
	}
	dbStructure.Users[usr.ID] = usr

	err = db.writeDB(dbStructure)
	if err != nil {
		return UserR{}, fmt.Errorf("%w", err)
	}

	var usrResp UserR

	usrResp.ID = usr.ID
	usrResp.Email = usr.Email

	return usrResp, nil
}
