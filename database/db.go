package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	model "github.com/Calgorr/Full-URL_Shortener/model"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	host     = "db"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "url_shortener"
)

var (
	db  *sql.DB
	err error
)

func RunMigrations() error {
	db, err := sql.Open("postgres", "postgres://"+user+":"+password+"@"+host+":"+strconv.Itoa(port)+"/"+"postgres"+"?sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
	CREATE DATABASE url_shortener;

	\\c url_shortener;

	CREATE TABLE IF NOT EXISTS users (
		userid SERIAL PRIMARY KEY,
		created_at DATE,
		username VARCHAR(255) UNIQUE,
		password VARCHAR(255)
	);

	CREATE TABLE IF NOT EXISTS url (
		ID SERIAL PRIMARY KEY,
		UserID INT,
		Longurl VARCHAR(255),
		shorturl VARCHAR(255) UNIQUE,
		used_times INT,
		created_at TIMESTAMP,
		last_used_at TIMESTAMP,
		FOREIGN KEY (UserID) REFERENCES users(userid)
	);

	CREATE OR REPLACE FUNCTION delete_expired_url() RETURNS TRIGGER AS $$
	BEGIN
		DELETE FROM url
		WHERE last_used_at <= NOW() - INTERVAL '1 day';

		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER trg_delete_expired_url
	AFTER INSERT OR UPDATE OR DELETE ON url
	FOR EACH ROW EXECUTE FUNCTION delete_expired_url();
	`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// Connect establishes a Connection to the database
func Connect() (*sql.DB, error) {
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
	Connect()
	defer db.Close()
	sqlStatement := "INSERT INTO users (created_at,username,password) VALUES ($1,$2,$3)"
	_, err := db.Exec(sqlStatement, time.Now(), user.Username, user.Password)
	return err
}

// GetUserByUsername retrieves a user from the database based on their username and password
func GetUserByUsername(user *model.User) (*model.User, error) {
	Connect()
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
	db, err := Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	if exists, _ := checkUserIDLongURL(int(id), link.LongURL); exists {
		return errors.New("URL already exists")
	}
	sqlstt := `INSERT INTO url (userid,longurl,shorturl,used_times,created_at,last_used_at) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err = db.Exec(sqlstt, int(id), link.LongURL, link.ShortURL, link.UsedTimes, link.CreatedAt, link.LastUsed)
	if err != nil {
		return err
	}
	return nil
}

func checkUserIDLongURL(userID int, longURL string) (bool, error) {
	db, err := Connect()
	if err != nil {
		return false, err
	}
	defer db.Close()
	sqlstt := `SELECT * FROM url WHERE userid=$1 AND longurl=$2`
	url := new(model.URL)
	err = db.QueryRow(sqlstt, userID, longURL).Scan(&url.ID, &url.UserID, &url.LongURL, &url.ShortURL, &url.UsedTimes, &url.CreatedAt, &url.LastUsed)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetLink retrieves a URL from the database based on its long URL
func GetLinkByLongURL(LongURL string, id float64) (*model.URL, error) {
	db, err := Connect()
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
	db, err := Connect()
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
	db, err := Connect()
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
	db, err := Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	sqlstt := `UPDATE url SET used_times=used_times+1, last_used_at=$1 WHERE shorturl=$2`
	_, err = db.Exec(sqlstt, time.Now(), shortURL)
	return err
}
