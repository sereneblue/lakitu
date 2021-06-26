package middleware

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func RequireLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)

		if sess.Values["authenticated"] == true {
			next(c)
			return nil
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
		})
	}
}
