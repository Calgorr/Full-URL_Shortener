package handler

import (
	"fmt"
	"net/http"

	authentication "github.com/Calgorr/IE_Backend_Fall/Authentication"
	model "github.com/Calgorr/IE_Backend_Fall/Model"
	"github.com/Calgorr/IE_Backend_Fall/database"
	"github.com/labstack/echo/v4"
)

func Signup(c echo.Context) error {
	newUser, err := new(model.User).Bind(c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, "Something went wrong")
	}
	err = database.AddUser(newUser)
	return c.String(http.StatusOK, "Signed up")
}

func Login(c echo.Context) error {
	user, err := new(model.User).Bind(c)
	if err != nil {
		return c.String(http.StatusBadRequest, "Something went wrong")
	}
	u, err := database.GetUserByUsername(user.Username)
	if err != nil || u.Password != user.Password {
		return c.String(http.StatusUnauthorized, "Wrong username or password")
	}
	id, err := database.GetIDByUsername(user.Username)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}
	token, err := authentication.GenerateJWT(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}
	c.Response().Header().Set("Authorization", token)
	return c.String(http.StatusOK, "Logged in")
}
