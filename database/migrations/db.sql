CREATE DATABASE url_shortener;

\c url_shortener;

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
    shorturl VARCHAR(255),
    used_times INT,
    created_at TIMESTAMP,
    last_used_at TIMESTAMP,
    FOREIGN KEY (UserID) REFERENCES users(userid)
    CONSTRAINT unique_user_Longurl UNIQUE (UserID, Longurl)
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
