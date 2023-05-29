package main

import (
	routes "github.com/Calgorr/Full-URL_Shortener/Routes"
	"github.com/labstack/echo/v4"
)

func main() {

	// Echo setup
	e := echo.New()

	// Routes
	e.POST("/urls", routes.CreateURL)
	e.GET("/:shortURL", RedirectURL)
	e.GET("/:shortURL/stats", GetURLStats)

	// Start the server
	e.Start(":8080")
}
