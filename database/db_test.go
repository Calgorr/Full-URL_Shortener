package database

import (
	"testing"

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
