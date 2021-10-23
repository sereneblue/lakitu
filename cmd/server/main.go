package main

import (
	"log"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/sereneblue/lakitu/internal/middleware"
	"github.com/sereneblue/lakitu/internal/routes"
	"github.com/sereneblue/lakitu/models"
)

var CookieStoreSecret string

func main() {
	err := models.IsInit()
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
	sess.POST("/change-password", routes.ChangePassword)
	sess.POST("/update-preferences", routes.ChangePreferences)
	sess.POST("/login", routes.Login)
	sess.GET("/logout", routes.Logout)
	sess.GET("/user", routes.UserData, middleware.RequireLogin)

	aws := e.Group("/aws")
	aws.POST("/verify", routes.VerifiyAWSCredentials)
	aws.POST("/ping", routes.PingAWS)
	aws.GET("/regions", routes.GetAWSRegions, middleware.RequireLogin)
	aws.POST("/gpu-instances", routes.GetAWSGPUInstances, middleware.RequireLogin)
	aws.POST("/pricing", routes.GetAWSPricing, middleware.RequireLogin)

	jobs := e.Group("/jobs")
	jobs.Use(middleware.RequireLogin)
	jobs.GET("", routes.GetCurrentJobStatus)
	jobs.GET("/:id", routes.GetJobStatus)

	machine := e.Group("/machine")
	machine.Use(middleware.RequireLogin)
	machine.GET("/list", routes.ListMachines)
	machine.POST("/create", routes.CreateMachine)
	machine.POST("/delete", routes.DeleteMachine)
	machine.POST("/start", routes.StartMachine)
	machine.POST("/stop", routes.StopMachine)
	machine.POST("/resize", routes.ResizeMachine)
	machine.POST("/transfer", routes.TransferMachine)

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
