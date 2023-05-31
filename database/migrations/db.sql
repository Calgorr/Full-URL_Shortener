-- Create the database if it doesn't exist
CREATE DATABASE url_shortener;

-- Switch to the newly created database
\ c url_shortener;

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
    shorturl VARCHAR(255) UNIQUE,
    used_times INT,
    created_at TIMESTAMP,
    last_used_at TIMESTAMP,
    FOREIGN KEY (UserID) REFERENCES users(userid)
);

-- Create the trigger function
CREATE
OR REPLACE FUNCTION delete_expired_url() RETURNS TRIGGER AS $ $ BEGIN
DELETE FROM
    url
WHERE
    last_used_at <= NOW() - INTERVAL '1 day';

RETURN NULL;

END;

$ $ LANGUAGE plpgsql;

-- Create the trigger
CREATE TRIGGER trg_delete_expired_url
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON url FOR EACH ROW EXECUTE FUNCTION delete_expired_url();