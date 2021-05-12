package main

import (
	"log"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/sereneblue/lakitu/internal/routes"
	"github.com/sereneblue/lakitu/models"
)

var CookieStoreSecret string

func main() {
	err := models.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	defer models.CloseDB()

	e := echo.New()
	e.HideBanner = true
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(CookieStoreSecret))))

	e.GET("/firstrun", routes.FirstRunCheck)
	e.POST("/verify-creds", routes.VerifiyCredentials)
	e.POST("/complete-setup", routes.CompleteSetup)
	e.POST("/ping-aws", routes.PingAWS)

	sess := e.Group("/session")
	sess.POST("/login", routes.Login)
	sess.GET("/logout", routes.Logout)

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
