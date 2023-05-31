package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	model "github.com/Calgorr/Full-URL_Shortener/model"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "calgor"
	password = "ami1r3ali"
	dbname   = "url_shortener"
)

var (
	db  *sql.DB
	err error
)

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
func AddUser(user *model.User) error {
	connect()
	defer db.Close()
	sqlStatement := "INSERT INTO users (created_at,username,password) VALUES ($1,$2,$3)"
	_, err := db.Exec(sqlStatement, time.Now(), user.Username, user.Password)
	fmt.Println(err)
	return err
}
func GetUserByUsername(user *model.User) (*model.User, error) {
	connect()
	defer db.Close()
	sqlStatement := "SELECT userid,username, password FROM users WHERE username=$1 AND password=$2"
	row := db.QueryRow(sqlStatement, user.Username, user.Password)
	u := new(model.User)
	err := row.Scan(&u.UserID, &u.Username, &u.Password)
	if err == sql.ErrNoRows {
		return nil, errors.New("User does not exists")
	}
	return u, nil
}

func AddLink(link *model.URL, id float64) error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlstt := `INSERT INTO url (userid,longurl,shorturl,used_times,created_at) VALUES ($1,$2,$3,$4,$5)`
	_, err = db.Exec(sqlstt, int(id), link.LongURL, link.ShortURL, link.UsedTimes, link.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func GetLink(shortURL string) (*model.URL, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sqlstt := `SELECT * FROM url WHERE shorturl=$1`
	url := new(model.URL)
	err = db.QueryRow(sqlstt, shortURL).Scan(&url.ID, &url.UserID, &url.LongURL, &url.ShortURL, &url.UsedTimes, &url.CreatedAt)
	if err != nil {
		return nil, err
	}
	return url, nil
}

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

func IncrementUsage(shortURL string) error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlstt := `UPDATE url SET used_times=used_times+1 WHERE shorturl=$1`
	_, err = db.Exec(sqlstt, shortURL)
	return err
}
