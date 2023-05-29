// db.go

package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// Database represents the database connection.
type Database struct {
	Conn *sql.DB
}

// InitializeDB initializes the database connection.
func InitializeDB() (*Database, error) {
	connStr := "postgres://calgor:ami1r3ali@localhost/url_shortener?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
		return nil, err
	}

	return &Database{Conn: db}, nil
}

// CreateUser creates a new user in the database.
func (db *Database) CreateUser(username, password string) error {
	stmt, err := db.Conn.Prepare("INSERT INTO users (username, password) VALUES ($1, $2)")
	if err != nil {
		log.Println("Failed to prepare statement:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, password)
	if err != nil {
		log.Println("Failed to execute query:", err)
		return err
	}

	return nil
}

// GetUserByUsername retrieves a user from the database by username.
func (db *Database) GetUserByUsername(username string) (*User, error) {
	var user User
	err := db.Conn.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		log.Println("Failed to execute query:", err)
		return nil, err
	}

	return &user, nil
}

// CreateURL creates a new URL in the database.
func (db *Database) CreateURL(userID int64, longURL string) (*URL, error) {
	shortURL := GenerateShortURL()

	stmt, err := db.Conn.Prepare("INSERT INTO urls (user_id, long_url, short_url) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		log.Println("Failed to prepare statement:", err)
		return nil, err
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(userID, longURL, shortURL).Scan(&id)
	if err != nil {
		log.Println("Failed to execute query:", err)
		return nil, err
	}

	return &URL{
		ID:       id,
		UserID:   userID,
		LongURL:  longURL,
		ShortURL: shortURL,
	}, nil
}

// GetURLByShortURL retrieves a URL from the database by its short URL.
func (db *Database) GetURLByShortURL(shortURL string) (*URL, error) {
	var url URL
	err := db.Conn.QueryRow("SELECT id, user_id, long_url, short_url FROM urls WHERE short_url = $1", shortURL).Scan(&url.ID, &url.UserID, &url.LongURL, &url.ShortURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // URL not found
		}
		log.Println("Failed to execute query:", err)
		return nil, err
	}

	return &url, nil
}

// GetUserURLs retrieves all URLs associated with a user from the database.
func (db *Database) GetUserURLs(userID int64) ([]URL, error) {
	rows, err := db.Conn.Query("SELECT id, user_id, long_url, short_url FROM urls WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Failed to execute query:", err)
		return nil, err
	}
	defer rows.Close()

	var urls []URL
	for rows.Next() {
		var url URL
		err = rows.Scan(&url.ID, &url.UserID, &url.LongURL, &url.ShortURL)
		if err != nil {
			log.Println("Failed to scan row:", err)
			continue
		}
		urls = append(urls, url)
	}

	return urls, nil
}

// GetURLStats retrieves the number of clicks/visits for a short URL.
func (db *Database) GetURLStats(shortURL string) (int, error) {
	var count int
	err := db.Conn.QueryRow("SELECT COUNT(*) FROM urls WHERE short_url = $1", shortURL).Scan(&count)
	if err != nil {
		log.Println("Failed to execute query:", err)
		return 0, err
	}

	return count, nil
}

// MigrateDB creates the necessary tables if they don't exist.
func MigrateDB(db *Database) error {
	_, err := db.Conn.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	)`)

	if err != nil {
		return err
	}

	_, err = db.Conn.Exec(`CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		long_url TEXT NOT NULL,
		short_url TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		FOREIGN KEY (user_id) REFERENCES users (id)
	)`)

	if err != nil {
		return err
	}

	return nil
}
