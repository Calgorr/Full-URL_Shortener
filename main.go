package main

import "github.com/labstack/echo/v4"

func main() {

	// Echo setup
	e := echo.New()

	// Routes
	e.POST("/urls", CreateURL)
	e.GET("/:shortURL", RedirectURL)
	e.GET("/:shortURL/stats", GetURLStats)

	// Start the server
	e.Start(":8080")
}
