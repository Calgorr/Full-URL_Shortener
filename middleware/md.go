package middleware

import (
	"github.com/labstack/echo/v4"
)

func TokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err == nil {
			c.Request().Header.Set(echo.HeaderAuthorization, cookie.Value)
		} else {
			return c.String(401, "unauthorized")
		}
		return next(c)
	}
}
