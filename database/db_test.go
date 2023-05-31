package database

import (
	"testing"
	"time"

	model "github.com/Calgorr/Full-URL_Shortener/model"
)

// TestAddUser tests the AddUser function
func TestAddUser(t *testing.T) {
	user := &model.User{
		Username: "testuser",
		Password: "testpassword",
	}
	err := AddUser(user)
	if err != nil {
		t.Errorf("AddUser failed: %v", err)
	}
}

// TestGetUserByUsername tests the GetUserByUsername function
func TestGetUserByUsername(t *testing.T) {
	user := &model.User{
		Username: "testuser",
		Password: "testpassword",
	}
	u, err := GetUserByUsername(user)
	if err != nil {
		t.Errorf("GetUserByUsername failed: %v", err)
	}
	// Verify that the retrieved user matches the expected user
	if u.Username != user.Username || u.Password != user.Password {
		t.Errorf("GetUserByUsername returned incorrect user")
	}
}

// TestAddLink tests the AddLink function
func TestAddLink(t *testing.T) {
	link := &model.URL{
		LongURL:   "http://example.com",
		ShortURL:  "abc123",
		UsedTimes: 0,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}
	err := AddLink(link, 1.0)
	if err != nil {
		t.Errorf("AddLink failed: %v", err)
	}
}
