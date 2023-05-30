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
	dbname   = "http_monitor"
)

var (
	db  *sql.DB
	err error
)

func connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	return db, err
}
func AddUser(user *model.User) error {
	fmt.Println(user, "moz")
	connect()
	defer db.Close()
	sqlStatement := "INSERT INTO users (created_at,username,password) VALUES ($1,$2,$3)"
	_, err := db.Exec(sqlStatement, time.Now(), user.Username, user.Password)
	return err
}
func GetUserByUsername(username string) (*model.User, error) {
	connect()
	defer db.Close()
	sqlStatement := "SELECT userid,username, password FROM users WHERE username=$1 "
	row := db.QueryRow(sqlStatement, username)
	user := new(model.User)
	err := row.Scan(&user.UserID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return nil, errors.New("User does not exists")
	}
	return user, nil
}

func AddLink(link *model.URL, id int) error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlstt := `INSERT INTO link (userid,longurl,shorturl,usedtimes,date) VALUES ($1,$2,$3,$4,$5)`
	_, err = db.Exec(sqlstt, id, link.LongURL, link.ShortURL, link.UsedTimes, link.CreatedAt)
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
	sqlstt := `SELECT * FROM link WHERE shorturl=$1`
	var link model.URL
	err = db.QueryRow(sqlstt, shortURL).Scan(&link.LongURL, &link.ShortURL, &link.UsedTimes)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func DeleteLink(shortURL string) error {
	db, err := connect()
	if err != nil {
		return errors.New("Internal Server Error")
	}
	defer db.Close()
	sqlstt := `DELETE FROM link WHERE shorturl=$1`
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
	sqlstt := `UPDATE link SET usedtimes=usedtimes+1 WHERE shorturl=$1`
	_, err = db.Exec(sqlstt, shortURL)
	return err
}
