package main

import (
    "log"

	"github.com/labstack/echo/v4"

    "github.com/sereneblue/lakitu/models"
    "github.com/sereneblue/lakitu/internal/routes"
)

func main() {
    err := models.InitDB()
    if err != nil {
        log.Fatal(err)
    }

    e := echo.New()
    e.HideBanner = true

    e.GET("/firstrun", routes.FirstRunCheck)

    e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}