package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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

func connect() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err := db.Ping()
	if err != nil {
		panic(err)
	}
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
	sqlStatement := "SELECT username, password FROM users WHERE username=$1 "
	row := db.QueryRow(sqlStatement, username)
	user := new(model.User)
	err := row.Scan(&user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return nil, errors.New("User does not exists")
	}
	return user, nil
}
func CreateURL(userID int64, longURL string) (*model.URL, error) {
	connect()
	defer db.Close()
	shortURL := model.GenerateShortURL()
	sqlState := "INSERT INTO urls (user_id, long_url, short_url) VALUES ($1, $2, $3) RETURNING id"

	var id int64
	err = db.QueryRow(sqlState, userID, longURL, shortURL).Scan(&id)
	if err != nil {
		log.Println("Failed to execute query:", err)
		return nil, err
	}

	return &model.URL{
		ID:        id,
		UserID:    userID,
		LongURL:   longURL,
		ShortURL:  shortURL,
		CreatedAt: time.Now(),
	}, nil
}
