package routes

import (
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/sereneblue/lakitu/internal/util"
	"github.com/sereneblue/lakitu/models"
)

func Login(c echo.Context) error {
	loginPwd := c.FormValue("password")

	match, err := argon2id.ComparePasswordAndHash(loginPwd, models.GetPasswordHash())
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "There was an error during login",
		})
	}

	if !match {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Invalid password",
		})
	}

	encKey, salt := models.GetEncryptedData()
	accessKey, secretKey := models.GetAWSSettings()

	// generate KEK
	KEK := util.CreateHashWithSalt(loginPwd, []byte(salt), models.KEKParams)

	key, err := util.Decrypt(encKey, KEK)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Could not decrypt encryption key",
		})
	}

	awsAccessKey, err := util.Decrypt(accessKey, key)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Could not decrypt AWS access key",
		})
	}

	awsSecretKey, err := util.Decrypt(secretKey, key)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": false,
			"message": "Could not decrypt AWS secret key",
		})
	}

	sess, _ := session.Get("session", c)
	sess.Values["authenticated"] = true
	sess.Values["accessKey"] = awsAccessKey
	sess.Values["secretKey"] = awsSecretKey

	sess.Save(c.Request(), c.Response())

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

func Logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}
