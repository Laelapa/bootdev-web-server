package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/Laelapa/bootdev-web-server/internal/authentication"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailInUse = errors.New("email already in use")
var ErrWrongCredentials = errors.New("wrong username or password")
var ErrInvalidRefToken = errors.New("invalid or expired refresh token")

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
	usrResp.IsChirpyRed = usr.IsChirpyRed

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
	userSan.RefToken, err = authentication.GenerateRefreshToken()
	if err != nil {
		return UserR{}, fmt.Errorf("error while trying to attribute a refresh token to the user: %w", err)
	}
	user.RefToken = userSan.RefToken
	user.RefExp = time.Now().Add(time.Hour * 24 * 60).Unix() // 60 days validity

	err = db.updateUserRefreshToken(user.ID, user.RefToken, user.RefExp)
	if err != nil {
		return UserR{}, fmt.Errorf("error while trying to register a refresh token to the db: %w", err)
	}

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

func (db *DB) UpdateUserCredentials(uID int, uEmail string, uPwd string) (User, error) {
	DBs, err := db.loadDB()
	if err != nil {
		return User{}, fmt.Errorf("%w", err)
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(uPwd), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("%w", err)
	}

	for i, v := range DBs.Users {
		if v.ID == uID {
			user := DBs.Users[i]
			user.Email = uEmail
			user.Password = pwdHash
			DBs.Users[i] = user

			err = db.writeDB(DBs)
			if err != nil {
				return User{}, fmt.Errorf("%w", err)
			}

			return User{ID: uID, Email: uEmail}, nil
		}
	}
	return User{}, ErrWrongCredentials
}

func (db *DB) updateUserRefreshToken(uID int, refToken string, refExp int64) error {
	dbData, err := db.loadDB()
	if err != nil {
		return fmt.Errorf("error while trying to pull the db: %w", err)
	}

	for i, v := range dbData.Users {
		if v.ID == uID {
			u := dbData.Users[i]
			u.RefToken = refToken
			u.RefExp = refExp
			dbData.Users[i] = u

			err = db.writeDB(dbData)
			if err != nil {
				return fmt.Errorf("error while trying to write the ref token to the db: %w", err)
			}

			return nil
		}
	}

	return ErrWrongCredentials
}

func (db *DB) CheckUserRefreshToken(refToken string) (userID int, err error) {
	dbData, err := db.loadDB()
	if err != nil {
		return 0, fmt.Errorf("error while trying to pull the db: %w", err)
	}

	for i := range dbData.Users {
		if refToken == dbData.Users[i].RefToken {
			if dbData.Users[i].RefExp > time.Now().Unix() {
				fmt.Printf("token checks out \n")
				return dbData.Users[i].ID, nil
			} else {
				return 0, ErrInvalidRefToken
			}
		}
	}

	return 0, ErrInvalidRefToken
}

func (db *DB) RevokeUserRefreshToken(refToken string) error {
	dbData, err := db.loadDB()
	if err != nil {
		return fmt.Errorf("error while trying to pull the db: %w", err)
	}

	for i := range dbData.Users {
		if refToken == dbData.Users[i].RefToken {
			user := dbData.Users[i]
			user.RefToken = ""
			user.RefExp = 0
			dbData.Users[i] = user

			err = db.writeDB(dbData)
			if err != nil {
				return fmt.Errorf("error while trying to del the ref token from the db: %w", err)
			}

			return nil
		}
	}

	return ErrInvalidRefToken
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
		ID:          len(users) + 1,
		Email:       email,
		Password:    pwdHash,
		IsChirpyRed: false,
	}
	dbStructure.Users[usr.ID] = usr

	err = db.writeDB(dbStructure)
	if err != nil {
		return UserR{}, fmt.Errorf("%w", err)
	}

	var usrResp UserR

	usrResp.ID = usr.ID
	usrResp.Email = usr.Email
	usrResp.IsChirpyRed = false

	return usrResp, nil
}

func (db *DB) UserP2W(uID int) (upgraded bool, err error) {
	dbData, err := db.loadDB()
	if err != nil {
		return false, fmt.Errorf("error while trying to pull the db: %w", err)
	}

	for i, v := range dbData.Users {
		if v.ID == uID {
			u := dbData.Users[i]
			u.IsChirpyRed = true
			dbData.Users[i] = u

			err = db.writeDB(dbData)
			if err != nil {
				return false, fmt.Errorf("error while trying to write the ref token to the db: %w", err)
			}

			return true, nil
		}
	}

	return false, nil
}
