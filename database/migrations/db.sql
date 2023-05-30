-- Create the database if it doesn't exist
CREATE DATABASE url_shortener;

-- Switch to the newly created database
\c url_shortener;

-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    userid SERIAL PRIMARY KEY,
    created_at DATE,
    username VARCHAR(255) UNIQUE,
    password VARCHAR(255)
);

-- Create the url table
CREATE TABLE IF NOT EXISTS url (
    ID SERIAL PRIMARY KEY,
    UserID INT,
    Longurl VARCHAR(255),
    shorturl VARCHAR(255),
    used_times INT,
    created_at DATE,
    FOREIGN KEY (UserID) REFERENCES users(userid)
);