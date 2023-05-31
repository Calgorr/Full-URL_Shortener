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
	fmt.Println(user)
	id, ok := userValidation(user)
	fmt.Println(id, ok)
	if !ok {
		return c.String(http.StatusUnauthorized, "invalid credentials")
	}
	token, err := authentication.GenerateJWT(id)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "internal server error")
	}
	c.Response().Header().Set(echo.HeaderAuthorization, token)
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: false,
	}
	c.SetCookie(cookie)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(user)
}

func userValidation(user *model.User) (int, bool) {
	user, err := database.GetUserByUsername(user)
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
	err = database.AddLink(link, claims["id"].(float64))
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "URL already exists")
	}
	return c.String(http.StatusOK, "Your Shortened link is "+c.Request().Host+"/"+link.ShortURL)
}

func Redirect(c echo.Context) error {
	if c.Param("shortURL") != "" {
		shortURL := c.Param("shortURL")
		address, err := database.GetLink(shortURL)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		if address.ShortURL != "" {
			database.IncrementUsage(address.ShortURL)
			return c.Redirect(http.StatusTemporaryRedirect, address.LongURL)
		} else {
			return c.String(http.StatusBadRequest, "Invalid url")
		}
	}
	return c.String(http.StatusInternalServerError, "Internal Server Error")
}

func GetURLStats(c echo.Context) error {
	shortURL := c.Param("shortURL")
	if shortURL == "" {
		return c.String(http.StatusBadRequest, "short url is required")
	}
	address, err := database.GetLink(shortURL)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	if address.ShortURL != "" {
		return c.JSON(http.StatusOK, "{\"shortURL\":\""+address.ShortURL+"\",\"longURL\":\""+address.LongURL+"\",\"usageCount\":"+fmt.Sprintf("%d", address.UsedTimes)+"}")
	} else {
		return c.String(http.StatusBadRequest, "Invalid url")
	}
}
