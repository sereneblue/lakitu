package main

import (
	"log"
	"mime"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/sereneblue/lakitu/assets"
	lakituMiddleware "github.com/sereneblue/lakitu/internal/middleware"
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

	mime.AddExtensionType(".js", "application/javascript")

	setup := e.Group("/setup")
	setup.GET("/first-run", routes.FirstRunCheck)
	setup.POST("/complete", routes.CompleteSetup)

	sess := e.Group("/session")
	sess.POST("/change-password", routes.ChangePassword, lakituMiddleware.RequireLogin)
	sess.POST("/update-preferences", routes.ChangePreferences, lakituMiddleware.RequireLogin)
	sess.POST("/login", routes.Login)
	sess.GET("/loggedin", routes.IsLoggedIn, lakituMiddleware.RequireLogin)
	sess.GET("/logout", routes.Logout, lakituMiddleware.RequireLogin)
	sess.GET("/user", routes.UserData, lakituMiddleware.RequireLogin)

	aws := e.Group("/aws")
	aws.POST("/verify", routes.VerifiyAWSCredentials)
	aws.POST("/ping", routes.PingAWS)
	aws.GET("/regions", routes.GetAWSRegions, lakituMiddleware.RequireLogin)
	aws.POST("/gpu-instances", routes.GetAWSGPUInstances, lakituMiddleware.RequireLogin)
	aws.POST("/pricing", routes.GetAWSPricing, lakituMiddleware.RequireLogin)

	jobs := e.Group("/jobs")
	jobs.Use(lakituMiddleware.RequireLogin)
	jobs.GET("", routes.GetCurrentJobStatus)
	jobs.GET("/:id", routes.GetJobStatus)

	machine := e.Group("/machine")
	machine.Use(lakituMiddleware.RequireLogin)
	machine.GET("/list", routes.ListMachines)
	machine.POST("/create", routes.CreateMachine)
	machine.POST("/delete", routes.DeleteMachine)
	machine.POST("/start", routes.StartMachine)
	machine.POST("/stop", routes.StopMachine)
	machine.POST("/resize", routes.ResizeMachine)
	machine.POST("/transfer", routes.TransferMachine)

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "build",
		Index:      "index.html",
		Browse:     false,
		HTML5:      true,
		Filesystem: http.FS(assets.App),
	}))

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
