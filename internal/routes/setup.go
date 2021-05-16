package routes

import (
	"net/http"

	"github.com/alexedwards/argon2id"

	"github.com/labstack/echo/v4"
	"github.com/sereneblue/lakitu/internal/util"
	"github.com/sereneblue/lakitu/models"
)

type SetupForm struct {
	DefaultRegion string `form:"region"`
	Password      string `form:"password"`
	AccessKey     string `form:"accessKey"`
	SecretKey     string `form:"secretKey"`
}

type LatencyForm struct {
	Regions []string `form:"regions"`
}

func CompleteSetup(c echo.Context) error {
	form := new(SetupForm)

	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Invalid input",
		})
	}

	var s models.Settings

	s.Key = "defaultRegion"
	s.Value = form.DefaultRegion
	s.Insert()

	// save password hash
	pwdHash, err := argon2id.CreateHash(form.Password, argon2id.DefaultParams)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Could not hash password",
		})
	}
	s.Key = "password"
	s.Value = pwdHash
	s.Insert()

	key, err := util.GenerateRandomKey()
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Could not generate random key",
		})
	}

	// generate KEK from password
	KEK, err := argon2id.CreateHash(form.Password, models.KEKParams)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Could not generate KEK",
		})
	}

	_, salt, _, _ := argon2id.DecodeHash(KEK)
	s.Key = "encSalt"
	s.Value = string(salt)
	s.Insert()

	// save encrypted key
	encKey, err := util.Encrypt(key, KEK)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Could not generate encrypted key",
		})
	}

	s.Key = "encKey"
	s.Value = encKey
	s.Insert()

	// save encrypted aws credentials
	encAWSAccessKey, err := util.Encrypt(form.AccessKey, key)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Could not encrypt AWS access key",
		})
	}
	s.Key = "awsAccessKeyId"
	s.Value = encAWSAccessKey
	s.Insert()

	encAWSSecretKey, err := util.Encrypt(form.SecretKey, key)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Could not encrypt AWS secret key",
		})
	}
	s.Key = "awsSecretKey"
	s.Value = encAWSSecretKey
	s.Insert()

	return c.JSON(http.StatusOK, map[string]bool{
		"success": true,
	})
}

func FirstRunCheck(c echo.Context) error {
	res := models.IsFirstRun()

	return c.JSON(http.StatusOK, map[string]bool{
		"success": res,
	})
}
