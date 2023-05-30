package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	authentication "github.com/Calgorr/Full-URL_Shortener/authentication"
	"github.com/Calgorr/Full-URL_Shortener/database"
	model "github.com/Calgorr/Full-URL_Shortener/model"
	"github.com/labstack/echo/v4"
)

func SignUp(c echo.Context) error {
	user := new(model.User)
	user, err := user.Bind(c)
	err = database.AddUser(user)
	if err != nil {
		return c.String(http.StatusConflict, "user already exists")
	}
	c.String(http.StatusOK, "success")
	return nil
}

func Login(c echo.Context) error {
	user := new(model.User)
	user, err := user.Bind(c)
	id, ok := userValidation(user)
	if !ok {
		return c.String(http.StatusUnauthorized, "invalid credentials")
	}
	token, err := authentication.GenerateJWT(id)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "internal server error")
	}
	c.Response().Header().Set(echo.HeaderAuthorization, token)

	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(user)
}

func userValidation(user *model.User) (int, bool) {
	user, err := database.GetUserByUsername(user.Username)
	if err != nil {
		return -1, false
	}
	return int(user.UserID), true
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
	token := c.Request().Header.Get("Authorization")
	claims, err := authentication.ExtractClaimsFromToken(token)
	if err != nil {
		return c.String(http.StatusUnauthorized, "invalid token")
	}
	link := model.NewLink(url)
	err = database.AddLink(link, claims["id"].(int))
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.String(http.StatusOK, "Your Shortened link is "+c.Request().Host+"/"+link.ShortURL)
}

func Redirect(c echo.Context) error {
	var err error
	var address *model.URL
	if c.Param("hash") != "" {
		hash := c.Param("hash")
		address, err = database.GetLink(hash)
		if address.ShortURL != "" {
			database.IncrementUsage(address.ShortURL)
			err = c.Redirect(http.StatusTemporaryRedirect, address.LongURL)
		} else {
			err = c.String(http.StatusBadRequest, "Invalid url")
		}
	}
	return err
}
