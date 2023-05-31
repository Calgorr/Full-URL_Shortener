package main

import (
	routes "github.com/Calgorr/Full-URL_Shortener/Routes"
	"github.com/Calgorr/Full-URL_Shortener/authentication"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo setup
	e := echo.New()
	e.Use(middleware.Logger())
	e.POST("/signup", routes.SignUp)
	e.GET("/signup", func(c echo.Context) error {
		return c.File("form/signup.html")
	})
	e.POST("/login", routes.Login)
	e.GET("/login", func(c echo.Context) error {
		return c.File("form/login.html")
	})
	// Routes
	v := e.Group("")
	v.Use(authentication.ValidateJWT)
	e.POST("/urls", routes.SaveUrl)
	e.GET("/:shortURL", routes.Redirect)
	//e.GET("/:shortURL/stats", GetURLStats)

	// Start the server
	e.Start(":8080")
}
