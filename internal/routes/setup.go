package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sereneblue/lakitu/models"
)

type AWSCredentials struct {
	AccessKey string `form:"accessKey"`
	SecretKey string `form:"secretKey"`
}

func FirstRunCheck(c echo.Context) error {
	res := models.IsFirstRun()

	return c.JSON(http.StatusOK, map[string]bool{
		"success": res,
	})
}

func VerifiyCredentials(c echo.Context) error {
	creds := new(AWSCredentials)

	if err := c.Bind(creds); err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Invalid input",
		})
	}

	client := models.NewAWSClient(creds.AccessKey, creds.SecretKey, "us-east-1")
	success, err := client.IsValidAWSCredentials()

	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": success,
			"message": err.Error(),
		})
	}

	regions := client.GetRegions()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": success,
		"data": map[string]interface{}{
			"regions": regions,
		},
	})
}
