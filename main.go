package main

import (
	routes "github.com/Calgorr/Full-URL_Shortener/Routes"
	"github.com/Calgorr/Full-URL_Shortener/authentication"
	md "github.com/Calgorr/Full-URL_Shortener/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// fmt.Println(database.RunMigrations())
	// res, err := database.Connect()
	// fmt.Println(res, err)
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
	v.Use(md.TokenMiddleware)
	v.Use(authentication.ValidateJWT)
	v.POST("/urls", routes.SaveUrl)
	e.GET("/urls", func(c echo.Context) error {
		return c.File("form/url.html")
	})
	v.GET("/:shortURL", routes.Redirect)
	v.GET("/:shortURL/stats", routes.GetURLStats)

	// Start the server
	e.Start(":8080")
}
