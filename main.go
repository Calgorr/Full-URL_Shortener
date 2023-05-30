package main

import (
	routes "github.com/Calgorr/Full-URL_Shortener/Routes"
	"github.com/Calgorr/Full-URL_Shortener/authentication"
	"github.com/labstack/echo/v4"
)

func main() {

	// Echo setup
	e := echo.New()

	e.POST("/signup", routes.SignUp)
	e.POST("login", routes.Login)
	// Routes
	v := e.Group("")
	v.Use(authentication.ValidateJWT)
	e.POST("/urls", routes.SaveUrl)
	e.GET("/:shortURL", routes.Redirect)
	//e.GET("/:shortURL/stats", GetURLStats)

	// Start the server
	e.Start(":8080")
}
