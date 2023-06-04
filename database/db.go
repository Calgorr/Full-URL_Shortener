package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	model "github.com/Calgorr/Full-URL_Shortener/model"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "ami1r3ali"
	dbname   = "url_shortener"
)

var (
	db  *sql.DB
	err error
)

// connect establishes a connection to the database
func connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	return db, err
}

// AddUser inserts a new user into the database
func AddUser(user *model.User) error {
	connect()
	defer db.Close()
	sqlStatement := "INSERT INTO users (created_at,username,password) VALUES ($1,$2,$3)"
	_, err := db.Exec(sqlStatement, time.Now(), user.Username, user.Password)
	return err
}

// GetUserByUsername retrieves a user from the database based on their username and password
func GetUserByUsername(user *model.User) (*model.User, error) {
	connect()
	defer db.Close()
	sqlStatement := "SELECT userid,username, password FROM users WHERE username=$1 AND password=$2"
	row := db.QueryRow(sqlStatement, user.Username, user.Password)
	u := new(model.User)
	err := row.Scan(&u.UserID, &u.Username, &u.Password)
	if err == sql.ErrNoRows {
		return nil, errors.New("User does not exist")
	}
	return u, nil
}

// AddLink inserts a new URL into the database
func AddLink(link *model.URL, id float64) error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlstt := `INSERT INTO url (userid,longurl,shorturl,used_times,created_at,last_used_at) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err = db.Exec(sqlstt, int(id), link.LongURL, link.ShortURL, link.UsedTimes, link.CreatedAt, link.LastUsed)
	if err != nil {
		return err
	}
	return nil
}

// GetLink retrieves a URL from the database based on its long URL
func GetLinkByLongURL(LongURL string, id float64) (*model.URL, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sqlstt := `SELECT * FROM url WHERE longurl=$1 AND userid=$2`
	url := new(model.URL)
	err = db.QueryRow(sqlstt, LongURL, id).Scan(&url.ID, &url.UserID, &url.LongURL, &url.ShortURL, &url.UsedTimes, &url.CreatedAt, &url.LastUsed)
	if err != nil {
		return nil, err
	}
	return url, nil
}

// GetLinkByShortURL retrieves a URL from the database based on its short URL
func GetLinkByShortURL(shortURL string) (*model.URL, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sqlstt := `SELECT * FROM url WHERE shorturl=$1`
	url := new(model.URL)
	err = db.QueryRow(sqlstt, shortURL).Scan(&url.ID, &url.UserID, &url.LongURL, &url.ShortURL, &url.UsedTimes, &url.CreatedAt, &url.LastUsed)
	if err != nil {
		return nil, err
	}
	return url, nil
}

// DeleteLink deletes a URL from the database based on its short URL
func DeleteLink(shortURL string) error {
	db, err := connect()
	if err != nil {
		return errors.New("Internal Server Error")
	}
	defer db.Close()
	sqlstt := `DELETE FROM url WHERE shorturl=$1`
	_, err = db.Exec(sqlstt, shortURL)
	if err != nil {
		return errors.New("Internal Server Error")
	}
	return nil
}

// IncrementUsage increments the usage count for a URL in the database
func IncrementUsage(shortURL string) error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlstt := `UPDATE url SET used_times=used_times+1, last_used_at=$1 WHERE shorturl=$2`
	_, err = db.Exec(sqlstt, time.Now(), shortURL)
	return err
}
