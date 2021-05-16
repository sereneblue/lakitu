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

	setup := e.Group("/setup")
	setup.GET("/first-run", routes.FirstRunCheck)
	setup.POST("/complete", routes.CompleteSetup)

	sess := e.Group("/session")
	sess.POST("/login", routes.Login)
	sess.GET("/logout", routes.Logout)

	aws := e.Group("/aws")
	aws.GET("/regions", routes.GetAWSRegions)
	aws.POST("/verify", routes.VerifiyAWSCredentials)
	aws.POST("/ping", routes.PingAWS)
	aws.POST("/gpu-instances", routes.GetAWSGPUInstances)

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
