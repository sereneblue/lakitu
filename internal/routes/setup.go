package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sereneblue/lakitu/models"
)

func FirstRunCheck(c echo.Context) error {
	res := models.IsFirstRun()

	return c.JSON(http.StatusOK, map[string]bool{
		"success": res,
	})
}