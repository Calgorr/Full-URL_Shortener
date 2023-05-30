package model

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

type User struct {
	UserID   int64  `form:"userid" json:"userid"`
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

func NewUser(username, password string) (*User, error) {
	if len(password) == 0 || len(username) == 0 {
		return nil, errors.New("username or password can not be empty")
	}
	return &User{Username: username, Password: password}, nil
}

func (u *User) Bind(c echo.Context) (*User, error) {
	fmt.Println(u)
	c.Bind(u)
	return u, nil
}
