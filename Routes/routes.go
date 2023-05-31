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

// SignUp handles the sign-up request
func SignUp(c echo.Context) error {
	user := new(model.User)
	user, err := user.Bind(c)    // Binding user data from the request
	err = database.AddUser(user) // Adding user to the database
	if err != nil {
		return c.String(http.StatusConflict, "user already exists")
	}
	c.String(http.StatusOK, "success")
	return nil
}

// Login handles the login request
func Login(c echo.Context) error {
	user := new(model.User)
	user, err := user.Bind(c)      // Binding user data from the request
	id, ok := userValidation(user) // Validating the user's credentials
	if !ok {
		return c.String(http.StatusUnauthorized, "invalid credentials")
	}
	token, err := authentication.GenerateJWT(id) // Generating a JWT token for authentication
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}
	c.Response().Header().Set(echo.HeaderAuthorization, token) // Setting the Authorization header in the response
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: false,
	}
	c.SetCookie(cookie) // Setting the token as a cookie in the response
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(user) // Encoding user data in the response
}

// userValidation checks if the user credentials are valid
func userValidation(user *model.User) (int, bool) {
	user, err := database.GetUserByUsername(user) // Retrieving user data from the database
	if err != nil {
		return -1, false
	}
	return int(user.UserID), true
}

// SaveUrl handles the request to save a URL and generate a short link
func SaveUrl(c echo.Context) error {
	url := c.FormValue("url")
	customPath := c.FormValue("customPath")
	if url == "" {
		return c.String(http.StatusBadRequest, "url is required")
	}
	url = strings.Replace(url, "www.", "", -1)
	if !strings.Contains(url, "http://") {
		url = "http://" + url
	}

	token := c.Request().Header.Get("Authorization")            // Retrieving the token from the Authorization header
	claims, err := authentication.ExtractClaimsFromToken(token) // Extracting claims (user ID) from the token
	if err != nil {
		return c.String(http.StatusUnauthorized, "invalid token")
	}
	link := model.NewLink(url, customPath) // Creating a new link object
	if !link.IsValidURL(url) {
		return c.String(http.StatusBadRequest, "invalid url")
	}
	err = database.AddLink(link, claims["id"].(float64)) // Adding the link to the database
	if err != nil {
		return c.String(http.StatusInternalServerError, "URL already exists")
	}
	return c.String(http.StatusOK, "Your Shortened link is "+c.Request().Host+"/"+link.ShortURL)
}

// Redirect handles the redirect request for a short URL
func Redirect(c echo.Context) error {
	if c.Param("shortURL") != "" {
		shortURL := c.Param("shortURL")
		address, err := database.GetLink(shortURL) // Retrieving the original long URL from the short URL
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		if address.ShortURL != "" {
			database.IncrementUsage(address.ShortURL)                        // Incrementing the usage count for the short URL
			return c.Redirect(http.StatusTemporaryRedirect, address.LongURL) // Redirecting to the original long URL
		} else {
			return c.String(http.StatusBadRequest, "Invalid url")
		}
	}
	return c.String(http.StatusInternalServerError, "Internal Server Error")
}

// GetURLStats retrieves the statistics of a short URL
func GetURLStats(c echo.Context) error {
	shortURL := c.Param("shortURL")
	if shortURL == "" {
		return c.String(http.StatusBadRequest, "short url is required")
	}
	address, err := database.GetLink(shortURL) // Retrieving the link data from the short URL
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	if address.ShortURL != "" {
		return c.JSON(http.StatusOK, "{\"shortURL\":\""+address.ShortURL+"\",\"longURL\":\""+address.LongURL+"\",\"usageCount\":"+fmt.Sprintf("%d", address.UsedTimes)+"}")
	} else {
		return c.String(http.StatusBadRequest, "Invalid url")
	}
}
