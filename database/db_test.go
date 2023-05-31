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
