package routes

import (
	"fmt"
	"net/http"
	"strings"

	authentication "github.com/Calgorr/Full-URL_Shortener/authentication"
	"github.com/Calgorr/Full-URL_Shortener/database"
	model "github.com/Calgorr/Full-URL_Shortener/model"
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
	_, err = database.GetUserByUsername(u.Username)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}
	token, err := authentication.GenerateJWT()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}
	c.Response().Header().Set("Authorization", token)
	return c.String(http.StatusOK, "Logged in")
}

func SaveUrl(c echo.Context) error {
	url := c.FormValue("url")
	if url == "" {
		return c.String(http.StatusBadRequest, "url is required")
	}
	url = strings.Replace(url, "www.", "", -1)
	if !strings.Contains(url, "http://") {
		url = "http://" + url
	}
	link := model.NewLink(url)
	err := db.AddLink(link)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.String(http.StatusOK, "Your Shortened link is "+c.Request().Host+"/"+link.Hash)
}

func Redirect(c echo.Context) error {
	var err error
	var address string
	if c.Param("hash") != "" {
		hash := c.Param("hash")
		address, err = db.GetLink(hash)
		if address != "" {
			db.IncrementUsage(hash)
			err = c.Redirect(http.StatusTemporaryRedirect, address)
		} else {
			err = c.String(http.StatusBadRequest, "Invalid url")
		}
	}
	return err
}
